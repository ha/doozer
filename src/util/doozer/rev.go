package main

import (
	"fmt"
)

func init() {
	cmds["rev"] = cmd{rev, "", "show current revision"}
	cmdHelp["rev"] = "Prints the current revision of the store.\n"
}

func rev() {
	c := dial()

	rev, err := c.Rev()
	if err != nil {
		bail(err)
	}

	fmt.Println(rev)
}
