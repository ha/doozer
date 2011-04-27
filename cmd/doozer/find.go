package main

import (
	"fmt"
	"github.com/ha/doozer"
)


func init() {
	cmds["find"] = cmd{find, "<glob>", "list files"}
	cmdHelp["find"] = `Prints the tree matching <glob>

Rules for <glob> pattern-matching:
 - '?' matches a single char in a single path component
 - '*' matches zero or more chars in a single path component
 - '**' matches zero or more chars in zero or more components
 - any other sequence matches itself

Prints a sequence of paths, one for each file/directory. Format of each record:

  <path> LF

Here, <path> is the file's path, and LF is an ASCII line-feed char.
`
}


func find(glob string) {
	c, err := doozer.Dial(*addr)
	if err != nil {
		bail(err)
	}

	if glob[len(glob)-1:] != "/" {
		glob = glob + "/"
	}

	if *rrev == -1 {
		*rrev, err = c.Rev()
		if err != nil {
			bail(err)
		}
	}

	info, err := c.Walk(glob+"**", *rrev, 0, -1)
	if err != nil {
		bail(err)
	}

	for _, ev := range info {
		fmt.Println(ev.Path)
	}
}
