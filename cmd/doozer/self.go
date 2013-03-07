package main

import (
	"os"
)


func init() {
	cmds["self"] = cmd{self, "", "identify a node"}
	cmdHelp["self"] = "Prints the node's ID.\n"
}


func self() {
	c := dial()

	id, err := c.Self()
	if err != nil {
		bail(err)
	}

	os.Stdout.Write(id)
}
