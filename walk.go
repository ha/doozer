package doozer

type Visitor interface {
	VisitDir(path string, f *FileInfo) bool
	VisitFile(path string, f *FileInfo)
}

// Walk walks the file tree in revision rev, rooted at root,
// analogously to Walk in package path/filepath.
func Walk(c *Conn, rev int64, root string, v Visitor, errors chan<- error) {
	f, err := c.Statinfo(rev, root)
	if err != nil {
		if errors != nil {
			errors <- err
		}
		return
	}
	walk(c, rev, root, f, v, errors)
}

func walk(c *Conn, r int64, path string, f *FileInfo, v Visitor, errors chan<- error) {
	if !f.IsDir {
		v.VisitFile(path, f)
		return
	}

	if !v.VisitDir(path, f) {
		return
	}

	list, err := c.Getdirinfo(path, r, 0, -1)
	if err != nil && errors != nil {
		errors <- err
	}

	if path != "/" {
		path += "/"
	}
	for i := range list {
		walk(c, r, path+list[i].Name, &list[i], v, errors)
	}
}
