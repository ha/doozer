package doozer

import (
	"path/filepath"
)

// Walk walks the file tree in revision rev, rooted at root,
// analogously to Walk in package path/filepath.
func Walk(c *Conn, rev int64, root string, v WalkFunc) error {
	f, err := c.Statinfo(rev, root)
	if err != nil {
		v(root, f, err)
		return err
	}
	return walk(c, rev, root, f, v)
}

func walk(c *Conn, r int64, path string, f *FileInfo, v WalkFunc) (err error) {
	verr := v(path, f, nil)
	if !f.IsDir || verr == filepath.SkipDir {
		return
	}
	if verr != nil {
		return err
	}

	list, err := c.Getdirinfo(path, r, 0, -1)
	if err != nil {
		verr = v(path, f, err)
		if verr != nil {
			return err
		}
	}

	if path != "/" {
		path += "/"
	}
	for i := range list {
		err = walk(c, r, path+list[i].Name, &list[i], v)
		if err != nil {
			break
		}
	}

	return
}
