package doozer

import (
	"path/filepath"
)

// Walk walks the file tree in revision rev, rooted at root,
// analogously to Walk in package path/filepath.
func Walk(c *Conn, rev int64, root string, v WalkFunc) {
	f, err := c.Statinfo(rev, root)
	if err != nil {
		v(root, f, err)
		return
	}
	walk(c, rev, root, f, v)
}

func walk(c *Conn, r int64, path string, f *FileInfo, v WalkFunc) {
	if !f.IsDir {
		v(path, f, nil)
		return
	}

	if v(path, f, nil) == filepath.SkipDir {
		return
	}

	list, err := c.Getdirinfo(path, r, 0, -1)
	if err != nil {
		v(path, f, err)
	}

	if path != "/" {
		path += "/"
	}
	for i := range list {
		walk(c, r, path+list[i].Name, &list[i], v)
	}
}
