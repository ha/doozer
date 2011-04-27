package main

import (
	"fmt"
	"github.com/ha/doozer"
	"os"
)


func init() {
	cmds["wait"] = cmd{wait, "<glob>", "wait for a change"}
	cmdHelp["wait"] = `Prints the next change to a file matching <glob>.

If flag -r is given, prints the next change on or after <rev>.

Rules for <glob> pattern-matching:
 - '?' matches a single char in a single path component
 - '*' matches zero or more chars in a single path component
 - '**' matches zero or more chars in zero or more components
 - any other sequence matches itself

Output is a sequence of records, one for each change. Format of each record:

  <path> <rev> <set|del> <len> LF <body> LF

Here, <path> is the file's path, <rev> is the revision, <len> is the number of
bytes in the body, <body> is the bytes of the body, and LF is an ASCII
line-feed char.

If a file is deleted, <rev> will be 0.
`
}


func wait(path string) {
	c, err := doozer.Dial(*addr)
	if err != nil {
		bail(err)
	}

	if *rrev == -1 {
		*rrev, err = c.Rev()
		if err != nil {
			bail(err)
		}
	}

	ev, err := c.Wait(path, *rrev)
	if err != nil {
		bail(err)
	}

	var sd string
	switch {
	case ev.IsSet():
		sd = "set"
	case ev.IsDel():
		sd = "del"
	}
	fmt.Println(ev.Path, ev.Rev, sd, len(ev.Body))
	os.Stdout.Write(ev.Body)
	fmt.Println()
}
