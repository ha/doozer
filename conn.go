package doozer

import (
	"code.google.com/p/goprotobuf/proto"
	"encoding/binary"
	"errors"
	"github.com/kr/pretty"
	"time"

	"io"
	"log"
	"math/rand"
	"net"
	"net/url"
	"strings"
)

var (
	uriPrefix = "doozer:?"
)

var (
	ErrInvalidUri = errors.New("invalid uri")
)

type txn struct {
	req  request
	resp *response
	err  error
	done chan bool
}

type Conn struct {
	addr    string
	conn    net.Conn
	send    chan *txn
	msg     chan []byte
	err     error
	stop    chan bool
	stopped chan bool
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Dial connects to a single doozer server.
func Dial(addr string) (*Conn, error) {
	return dial(addr, -1)
}

// DialTimeout acts like Dial but takes a timeout.
func DialTimeout(addr string, timeout time.Duration) (*Conn, error) {
	return dial(addr, timeout)
}

func dial(addr string, timeout time.Duration) (*Conn, error) {
	var c Conn
	var err error
	c.addr = addr
	if timeout > 0 {
		c.conn, err = net.DialTimeout("tcp", addr, timeout)
	} else {
		c.conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		return nil, err
	}

	c.send = make(chan *txn)
	c.msg = make(chan []byte)
	c.stop = make(chan bool, 1)
	c.stopped = make(chan bool)
	errch := make(chan error, 1)
	go c.mux(errch)
	go c.readAll(errch)
	return &c, nil
}

// DialUri connects to one of the doozer servers given in `uri`. If `uri`
// contains a cluster name, it will lookup addrs to try in `buri`.  If `uri`
// contains a  secret key, then DialUri will call `Access` with the secret.
func DialUri(uri, buri string) (*Conn, error) {
	return dialUri(uri, buri, -1)
}

// DialUriTimeout acts like DialUri but takes a timeout.
func DialUriTimeout(uri, buri string, timeout time.Duration) (*Conn, error) {
	return dialUri(uri, buri, -1)
}

func dialUri(uri, buri string, timeout time.Duration) (*Conn, error) {
	if !strings.HasPrefix(uri, uriPrefix) {
		return nil, ErrInvalidUri
	}

	q := uri[len(uriPrefix):]
	p, err := url.ParseQuery(q)
	if err != nil {
		return nil, err
	}

	addrs := make([]string, 0)

	name, ok := p["cn"]
	if ok && buri != "" {
		c, err := DialUri(buri, "")
		if err != nil {
			return nil, err
		}

		addrs, err = lookup(c, name[0])
		if err != nil {
			return nil, err
		}
	} else {
		var ok bool
		addrs, ok = p["ca"]
		if !ok {
			return nil, ErrInvalidUri
		}
	}

	c, err := dial(addrs[rand.Int()%len(addrs)], timeout)
	if err != nil {
		return nil, err
	}

	secret, ok := p["sk"]
	if ok {
		err = c.Access(secret[0])
		if err != nil {
			c.Close()
			return nil, err
		}
	}

	return c, nil
}

// Find possible addresses for cluster named name.
func lookup(b *Conn, name string) (as []string, err error) {
	rev, err := b.Rev()
	if err != nil {
		return nil, err
	}

	path := "/ctl/ns/" + name
	names, err := b.Getdir(path, rev, 0, -1)
	if err, ok := err.(*Error); ok && err.Err == ErrNoEnt {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	path += "/"
	for _, name := range names {
		body, _, err := b.Get(path+name, &rev)
		if err != nil {
			return nil, err
		}
		as = append(as, string(body))
	}
	return as, nil
}

func (c *Conn) call(t *txn) error {
	t.done = make(chan bool)
	select {
	case <-c.stopped:
		return c.err
	case c.send <- t:
		select {
		case <-c.stopped:
			return c.err
		case <-t.done:
			if t.err != nil {
				return t.err
			}
			if t.resp.ErrCode != nil {
				return newError(t)
			}
		}
	}
	return nil
}

// After Close is called, operations on c will return ErrClosed.
func (c *Conn) Close() {
	select {
	case c.stop <- true:
	default:
	}
}

func (c *Conn) mux(errch chan error) {
	txns := make(map[int32]*txn)
	var n int32 // next tag
	var err error

	for {
		select {
		case t := <-c.send:
			// find an unused tag
			for t := txns[n]; t != nil; t = txns[n] {
				n++
			}
			txns[n] = t

			// don't take n's address; it will change
			tag := n
			t.req.Tag = &tag

			var buf []byte
			buf, err = proto.Marshal(&t.req)
			if err != nil {
				txns[n] = nil
				t.err = err
				t.done <- true
				continue
			}

			err = c.write(buf)
			if err != nil {
				goto error
			}
		case buf := <-c.msg:
			var r response
			err = proto.Unmarshal(buf, &r)
			if err != nil {
				log.Print(err)
				continue
			}

			if r.Tag == nil {
				log.Printf("nil tag: %# v", pretty.Formatter(r))
				continue
			}
			t := txns[*r.Tag]
			if t == nil {
				log.Printf("unexpected: %# v", pretty.Formatter(r))
				continue
			}

			delete(txns, *r.Tag)
			t.resp = &r
			t.done <- true
		case err = <-errch:
			goto error
		case <-c.stop:
			err = ErrClosed
			goto error
		}
	}

error:
	c.err = err
	for _, t := range txns {
		t.err = err
		t.done <- true
	}
	c.conn.Close()
	close(c.stopped)
}

func (c *Conn) readAll(errch chan error) {
	for {
		buf, err := c.read()
		if err != nil {
			errch <- err
			return
		}

		c.msg <- buf
	}
}

func (c *Conn) read() ([]byte, error) {
	var size int32
	err := binary.Read(c.conn, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(c.conn, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *Conn) write(buf []byte) error {
	err := binary.Write(c.conn, binary.BigEndian, int32(len(buf)))
	if err != nil {
		return err
	}

	_, err = c.conn.Write(buf)
	return err
}

// Attempts access to the store
func (c *Conn) Access(token string) error {
	var t txn
	t.req.Verb = request_ACCESS.Enum()
	t.req.Value = []byte(token)
	return c.call(&t)
}

// Sets the contents of file to body, if it hasn't been modified since oldRev.
func (c *Conn) Set(file string, oldRev int64, body []byte) (newRev int64, err error) {
	var t txn
	t.req.Verb = request_SET.Enum()
	t.req.Path = &file
	t.req.Value = body
	t.req.Rev = &oldRev

	err = c.call(&t)
	if err != nil {
		return
	}

	return t.resp.GetRev(), nil
}

// Deletes file, if it hasn't been modified since rev.
func (c *Conn) Del(file string, rev int64) error {
	var t txn
	t.req.Verb = request_DEL.Enum()
	t.req.Path = &file
	t.req.Rev = &rev
	return c.call(&t)
}

func (c *Conn) Nop() error {
	var t txn
	t.req.Verb = request_NOP.Enum()
	return c.call(&t)
}

// Returns the body and revision of the file at path,
// as of store revision *rev.
// If rev is nil, uses the current state.
func (c *Conn) Get(file string, rev *int64) ([]byte, int64, error) {
	var t txn
	t.req.Verb = request_GET.Enum()
	t.req.Path = &file
	t.req.Rev = rev

	err := c.call(&t)
	if err != nil {
		return nil, 0, err
	}

	return t.resp.Value, t.resp.GetRev(), nil
}

// Getdir reads up to lim names from dir, at revision rev, into an array.
// Names are read in lexicographical order, starting at position off.
// A negative lim means to read until the end.
func (c *Conn) Getdir(dir string, rev int64, off, lim int) (names []string, err error) {
	for lim != 0 {
		var t txn
		t.req.Verb = request_GETDIR.Enum()
		t.req.Rev = &rev
		t.req.Path = &dir
		t.req.Offset = proto.Int32(int32(off))
		err = c.call(&t)
		if err, ok := err.(*Error); ok && err.Err == ErrRange {
			return names, nil
		}
		if err != nil {
			return nil, err
		}
		names = append(names, *t.resp.Path)
		off++
		lim--
	}
	return
}

// Getdirinfo reads metadata for up to lim files from dir, at revision rev,
// into an array.
// Files are read in lexicographical order, starting at position off.
// A negative lim means to read until the end.
// Getdirinfo returns the array and an error, if any.
func (c *Conn) Getdirinfo(dir string, rev int64, off, lim int) (a []FileInfo, err error) {
	names, err := c.Getdir(dir, rev, off, lim)
	if err != nil {
		return nil, err
	}

	if dir != "/" {
		dir += "/"
	}
	a = make([]FileInfo, len(names))
	for i, name := range names {
		var fp *FileInfo
		fp, err = c.Statinfo(rev, dir+name)
		if err != nil {
			a[i].Name = name
		} else {
			a[i] = *fp
		}
	}
	return
}

// Statinfo returns metadata about the file or directory at path,
// in revision *storeRev. If storeRev is nil, uses the current
// revision.
func (c *Conn) Statinfo(rev int64, path string) (f *FileInfo, err error) {
	f = new(FileInfo)
	f.Len, f.Rev, err = c.Stat(path, &rev)
	if err != nil {
		return nil, err
	}
	if f.Rev == missing {
		return nil, ErrNoEnt
	}
	f.Name = basename(path)
	f.IsSet = true
	f.IsDir = f.Rev == dir
	return f, nil
}

// Stat returns metadata about the file or directory at path,
// in revision *storeRev. If storeRev is nil, uses the current
// revision.
func (c *Conn) Stat(path string, storeRev *int64) (len int, fileRev int64, err error) {
	var t txn
	t.req.Verb = request_STAT.Enum()
	t.req.Path = &path
	t.req.Rev = storeRev

	err = c.call(&t)
	if err != nil {
		return 0, 0, err
	}

	return int(t.resp.GetLen()), t.resp.GetRev(), nil
}

// Walk reads up to lim entries matching glob, in revision rev, into an array.
// Entries are read in lexicographical order, starting at position off.
// A negative lim means to read until the end.
// Conn.Walk will be removed in a future release. Use Walk instead.
func (c *Conn) Walk(glob string, rev int64, off, lim int) (info []Event, err error) {
	for lim != 0 {
		var t txn
		t.req.Verb = request_WALK.Enum()
		t.req.Rev = &rev
		t.req.Path = &glob
		t.req.Offset = proto.Int32(int32(off))
		err = c.call(&t)
		if err, ok := err.(*Error); ok && err.Err == ErrRange {
			return info, nil
		}
		if err != nil {
			return nil, err
		}
		info = append(info, Event{
			*t.resp.Rev,
			*t.resp.Path,
			t.resp.Value,
			*t.resp.Flags,
		})
		off++
		lim--
	}
	return
}

// Waits for the first change, on or after rev, to any file matching glob.
func (c *Conn) Wait(glob string, rev int64) (ev Event, err error) {
	var t txn
	t.req.Verb = request_WAIT.Enum()
	t.req.Path = &glob
	t.req.Rev = &rev

	err = c.call(&t)
	if err != nil {
		return
	}

	ev.Rev = *t.resp.Rev
	ev.Path = *t.resp.Path
	ev.Body = t.resp.Value
	ev.Flag = *t.resp.Flags & (set | del)
	return
}

// Waits for the first change, on or after rev, to any file matching glob,
// within the specific time expressed as time.Duration
func (c *Conn) WaitTimeout(glob string, rev int64, timeout time.Duration) (ev Event, err error) {
	var timer *time.Timer
	if timeout > 0 {
		timer = time.AfterFunc(timeout, func() {
			c.err = ErrWaitTimeout
			c.stopped <- true
		})
	}
	ev, err = c.Wait(glob, rev)
	if timer != nil {
		timer.Stop()
	}
	return
}

// Rev returns the current revision of the store.
func (c *Conn) Rev() (int64, error) {
	var t txn
	t.req.Verb = request_REV.Enum()

	err := c.call(&t)
	if err != nil {
		return 0, err
	}

	return *t.resp.Rev, nil
}

// Self returns the node's identifier
func (c *Conn) Self() ([]byte, error) {
	var t txn
	t.req.Verb = request_SELF.Enum()

	err := c.call(&t)
	if err != nil {
		return nil, err
	}

	return t.resp.Value, nil
}
