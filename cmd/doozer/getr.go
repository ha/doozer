package main

import (
	"fmt"
	"os"
)

func init() {
	cmds["getr"] = cmd{getr, "<path>", "read a file and its rev"}
	cmdHelp["getr"] = `Prints the body and rev of the file at <path>.

The format is:

  <path> <rev> <len> LF <body> LF

where LF is an ASCII line feed character.

`
}

func getr(path string) {
	c := dial()

	body, rev, err := c.Get(path, nil)
	if err != nil {
		bail(err)
	}

	fmt.Println(path, rev, len(body))
	os.Stdout.Write(body)
	fmt.Println()
}
