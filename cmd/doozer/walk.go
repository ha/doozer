package main

import (
	"fmt"
	"os"
)


func init() {
	cmds["walk"] = cmd{walk, "<glob>", "read many files"}
	cmdHelp["walk"] = `Prints the path, revision, and body of each file matching <glob>.

Rules for <glob> pattern-matching:
 - '?' matches a single char in a single path component
 - '*' matches zero or more chars in a single path component
 - '**' matches zero or more chars in zero or more components
 - any other sequence matches itself

Prints a sequence of records, one for each file. Format of each record:

  <path> <rev> <len> LF <body> LF

Here, <path> is the file's path, <rev> is the revision, <len> is the number of
bytes in the body, <body> is the bytes of the body, and LF is an ASCII
line-feed char.
`
}


func walk(glob string) {
	c := dial()

	if *rrev == -1 {
		var err os.Error
		*rrev, err = c.Rev()
		if err != nil {
			bail(err)
		}
	}

	info, err := c.Walk(glob, *rrev, 0, -1)
	if err != nil {
		bail(err)
	}

	for _, ev := range info {
		fmt.Println(ev.Path, ev.Rev, len(ev.Body))
		os.Stdout.Write(ev.Body)
		fmt.Println()
	}
}
