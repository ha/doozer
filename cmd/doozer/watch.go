package main

import (
	"fmt"
	"os"
)


func init() {
	cmds["watch"] = cmd{watch, "<glob>", "watch for a change"}
	cmdHelp["watch"] = `Prints changes to any file matching <glob>.

If flag -r is given, prints changes beginning on <rev>.

Rules for <glob> pattern-matching:
 - '?' matches a single char in a single path component
 - '*' matches zero or more chars in a single path component
 - '**' matches zero or more chars in zero or more components
 - any other sequence matches itself

Output is a sequence of records, one for each change. Format of each record:

  <path> <rev> <set|del> <len> LF <body> LF

Here, <path> is the file's path, <rev> is the revision of the change,
<len> is the number of bytes in the body, and LF is an ASCII line-feed char.
`
}


func watch(glob string) {
	c := dial()

	if *rrev == -1 {
		var err os.Error
		*rrev, err = c.Rev()
		if err != nil {
			bail(err)
		}
	}

	for {
		ev, err := c.Wait(glob, *rrev)
		if err != nil {
			bail(err)
		}
		*rrev = ev.Rev + 1

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
}
