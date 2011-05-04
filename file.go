package doozer

const (
	missing = int64(-iota)
	clobber
	dir
	nop
)

type FileInfo struct {
	Name  string
	Len   int
	Rev   int64
	IsSet bool
	IsDir bool
}

func basename(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
