package main

import (
	"fmt"
	"github.com/ha/doozer"
)


func init() {
	cmds["rev"] = cmd{rev, "<path>", "read a file"}
	cmdHelp["rev"] = "Prints the current revision.\n"
}


func rev() {
	c, err := doozer.Dial(*addr)
	if err != nil {
		bail(err)
	}

	rev, err := c.Rev()
	if err != nil {
		bail(err)
	}

	fmt.Println(rev)
}
