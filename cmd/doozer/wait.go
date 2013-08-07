package main

import (
	"flag"
	"fmt"
	"os"
)

var timeout = flag.Duration("t", 0, "wait timeout")

func init() {
	flag.Parse()
	cmds["wait"] = cmd{wait, "<glob>", "wait for a change"}
	cmdHelp["wait"] = `Prints the next change to a file matching <glob>.

If flag -r is given, prints the next change on or after <rev>.
if flag -t is given, the wait will timeout after the specified
	duration (i.e. -t 1m30s will wait 1 minute and 30 seconds,
	valid time units: "ns", "us" (or "µs"), "ms", "s", "m", "h").

Rules for <glob> pattern-matching:
 - '?' matches a single char in a single path component
 - '*' matches zero or more chars in a single path component
 - '**' matches zero or more chars in zero or more components
 - any other sequence matches itself

Output is a sequence of records, one for each change. Format of each record:

  <path> <rev> <set|del> <len> LF <body> LF

Here, <path> is the file's path, <rev> is the revision, <len> is the number of
bytes in the body, and LF is an ASCII line-feed char.
`
}

func wait(path string) {
	c := dial()

	if *rrev == -1 {
		var err error
		*rrev, err = c.Rev()
		if err != nil {
			bail(err)
		}
	}

	ev, err := c.WaitTimeout(path, *rrev, *timeout)
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
