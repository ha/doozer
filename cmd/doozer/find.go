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

	v := func(path string, f *doozer.FileInfo, e error) (err error) {
		if e != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Println(path)

		return nil
	}
	doozer.Walk(c, *rrev, path, v)
}
