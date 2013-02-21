package main

import (
	"os"
)

func init() {
	cmds["get"] = cmd{get, "<path>", "read a file"}
	cmdHelp["get"] = "Prints the body of the file at <path>.\n"
}

func get(path string) {
	c := dial()

	body, _, err := c.Get(path, nil)
	if err != nil {
		bail(err)
	}

	os.Stdout.Write(body)
}
