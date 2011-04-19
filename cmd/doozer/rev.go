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
	c := doozer.New("<test>", *addr)

	rev, err := c.Rev()
	if err != nil {
		bail(err)
	}

	fmt.Println(rev)
}
