package main

import (
	"fmt"
	"github.com/ha/doozer"
	"os"
)

func init() {
	cmds["find"] = cmd{find, "<path>", "list files"}
	cmdHelp["find"] = `Prints the tree rooted at <path>

Prints the path for each file or directory, one per line.
`
}

func find(path string) {
	c := dial()

	if *rrev == -1 {
		var err error
		*rrev, err = c.Rev()
		if err != nil {
			bail(err)
		}
	}

	v := make(vis)
	errs := make(chan error)
	go func() {
		doozer.Walk(c, *rrev, path, v, errs)
		close(v)
	}()

	for {
		select {
		case path, ok := <-v:
			if !ok {
				return
			}
			fmt.Println(path)
		case err := <-errs:
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

type vis chan string

func (v vis) VisitDir(path string, f *doozer.FileInfo) bool {
	v <- path
	return true
}

func (v vis) VisitFile(path string, f *doozer.FileInfo) {
	v <- path
}
