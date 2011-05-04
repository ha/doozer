package doozer


import (
	"encoding/binary"
	"github.com/kr/pretty.go"
	"goprotobuf.googlecode.com/hg/proto"
	"http"
	"io"
	"log"
	"net"
	"os"
	"rand"
	"strings"
)


var (
	uriPrefix = "doozer:?"
)

var (
	ErrInvalidUri = os.NewError("invalid uri")
)


type txn struct {
	req  request
	resp *response
	err  os.Error
	done chan bool
}


type Conn struct {
	addr    string
	conn    net.Conn
	send    chan *txn
	msg     chan []byte
	err     os.Error
	stop    chan bool
	stopped chan bool
}


// Dial connects to a single doozer server.
func Dial(addr string) (*Conn, os.Error) {
	var c Conn
	var err os.Error
	c.addr = addr
	c.conn, err = net.Dial("tcp", "", addr)
	if err != nil {
		return nil, err
	}

	c.send = make(chan *txn)
	c.msg = make(chan []byte)
	c.stop = make(chan bool, 1)
	c.stopped = make(chan bool)
	errch := make(chan os.Error, 1)
	go c.mux(errch)
	go c.readAll(errch)
	return &c, nil
}


func DialUri(uri string) (*Conn, os.Error) {
	if !strings.HasPrefix(uri, uriPrefix) {
		return nil, ErrInvalidUri
	}

	q := uri[len(uriPrefix):]
	p, err := http.ParseQuery(q)
	if err != nil {
		return nil, err
	}

	addrs, ok := p["ca"]
	if !ok {
		return nil, ErrInvalidUri
	}

	c := Dial(addrs[rand.Int()%len(addrs)])
}


func (c *Conn) call(t *txn) os.Error {
	t.done = make(chan bool)
	select {
	case <-c.stopped:
		return c.err
	case c.send <- t:
		<-t.done
		if t.err != nil {
			return t.err
		}
		if t.resp.ErrCode != nil {
			return newError(t)
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


func (c *Conn) mux(errch chan os.Error) {
	txns := make(map[int32]*txn)
	var n int32 // next tag
	var err os.Error

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

			txns[*r.Tag] = nil, false
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


func (c *Conn) readAll(errch chan os.Error) {
	for {
		buf, err := c.read()
		if err != nil {
			errch <- err
			return
		}

		c.msg <- buf
	}
}


func (c *Conn) read() ([]byte, os.Error) {
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


func (c *Conn) write(buf []byte) os.Error {
	err := binary.Write(c.conn, binary.BigEndian, int32(len(buf)))
	if err != nil {
		return err
	}

	_, err = c.conn.Write(buf)
	return err
}

// Attempts access to the store
func (c *Conn) Access(token string) os.Error {
	var t txn
	t.req.Verb = newRequest_Verb(request_ACCESS)
	t.req.Value = []byte(token)
	return c.call(&t)
}

// Sets the contents of file to body, if it hasn't been modified since oldRev.
func (c *Conn) Set(file string, oldRev int64, body []byte) (newRev int64, err os.Error) {
	var t txn
	t.req.Verb = newRequest_Verb(request_SET)
	t.req.Path = &file
	t.req.Value = body
	t.req.Rev = &oldRev

	err = c.call(&t)
	if err != nil {
		return
	}

	return proto.GetInt64(t.resp.Rev), nil
}


// Deletes file, if it hasn't been modified since rev.
func (c *Conn) Del(file string, rev int64) os.Error {
	var t txn
	t.req.Verb = newRequest_Verb(request_DEL)
	t.req.Path = &file
	t.req.Rev = &rev
	return c.call(&t)
}


func (c *Conn) Nop() os.Error {
	var t txn
	t.req.Verb = newRequest_Verb(request_NOP)
	return c.call(&t)
}


// Returns the body and revision of the file at path,
// as of store revision *rev.
// If rev is nil, uses the current state.
func (c *Conn) Get(file string, rev *int64) ([]byte, int64, os.Error) {
	var t txn
	t.req.Verb = newRequest_Verb(request_GET)
	t.req.Path = &file
	t.req.Rev = rev

	err := c.call(&t)
	if err != nil {
		return nil, 0, err
	}

	return t.resp.Value, proto.GetInt64(t.resp.Rev), nil
}


// Getdir reads up to lim names from dir, at revision rev, into an array.
// Names are read in lexicographical order, starting at position off.
// A negative lim means to read until the end.
func (c *Conn) Getdir(dir string, rev int64, off, lim int) (names []string, err os.Error) {
	for lim != 0 {
		var t txn
		t.req.Verb = newRequest_Verb(request_GETDIR)
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


// Stat returns metadata about the file or directory at path,
// in revision *storeRev. If storeRev is nil, uses the current
// revision.
func (c *Conn) Stat(path string, storeRev *int64) (len int32, fileRev int64, err os.Error) {
	var t txn
	t.req.Verb = newRequest_Verb(request_STAT)
	t.req.Path = &path
	t.req.Rev = storeRev

	err = c.call(&t)
	if err != nil {
		return 0, 0, err
	}

	return proto.GetInt32(t.resp.Len), proto.GetInt64(t.resp.Rev), nil
}


// Walk reads up to lim entries matching glob, in revision rev, into an array.
// Entries are read in lexicographical order, starting at position off.
// A negative lim means to read until the end.
func (c *Conn) Walk(glob string, rev int64, off, lim int) (info []Event, err os.Error) {
	for lim != 0 {
		var t txn
		t.req.Verb = newRequest_Verb(request_WALK)
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
func (c *Conn) Wait(glob string, rev int64) (ev Event, err os.Error) {
	var t txn
	t.req.Verb = newRequest_Verb(request_WAIT)
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


// Rev returns the current revision of the store.
func (c *Conn) Rev() (int64, os.Error) {
	var t txn
	t.req.Verb = newRequest_Verb(request_REV)

	err := c.call(&t)
	if err != nil {
		return 0, err
	}

	return *t.resp.Rev, nil
}
