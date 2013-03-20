package main

import (
	"os"
	"strings"
)

func init() {
	cmds["getdir"] = cmd{getdir, "<path>", "read a path"}
	cmdHelp["getdir"] = "Lists the contents of <path>.\n"
}

func getdir(path string) {
	c := dial()

	body, err := c.Getdir(path, nil, 0, -1)
	if err != nil {
		bail(err)
	}

	out := strings.Join(body, "\n")

	os.Stdout.Write([]byte(out))
}
