package main

import (
	"os"
	"strings"
)

func init() {
	cmds["ls"] = cmd{ls, "<dir>", "list files  in a directory"}
	cmdHelp["ls"] = "Prints the list of files under <dir>.\n"
}

func ls(path string) {
	c := dial()

	if *rrev == -1 {
		var err error
		*rrev, err = c.Rev()
		if err != nil {
			bail(err)
		}
	}

	body, err := c.Getdir(path, *rrev, 0, -1)
	if err != nil {
		bail(err)
	}

	os.Stdout.Write([]byte(strings.Join(body, "\n") + "\n"))
}
