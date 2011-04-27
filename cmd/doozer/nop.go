package main

import (
	"github.com/ha/doozer"
)


func init() {
	cmds["nop"] = cmd{nop, "", "consensus"}
	cmdHelp["nop"] = `Performs a consensus operation.

No change will be made to the data store.
`
}


func nop() {
	c, err := doozer.Dial(*addr)
	if err != nil {
		bail(err)
	}

	err = c.Nop()
	if err != nil {
		bail(err)
	}
}
