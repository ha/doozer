package main

import (
	"fmt"
	"os"
)


func init() {
	cmds["stat"] = cmd{stat, "<path>", "print file status"}
	cmdHelp["stat"] = `Print the status for each <path>.

If path is a directory, prints "d" and the number of entries.
Otherwise, prints its revision and length.
`
}


func stat(path string) {
	c := dial()

	len, rev, err := c.Stat(path, nil)
	if err != nil {
		bail(err)
	}

	switch rev {
	case 0:
		fmt.Fprintln(os.Stderr, "No such file or directory:", path)
		os.Exit(2)
	case -2:
		fmt.Println("d", len)
	default:
		fmt.Println(rev, len)
	}
}
