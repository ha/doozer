package main

import (
	"github.com/ha/doozer"
)


func init() {
	cmds["del"] = cmd{del, "<path> <rev>", "delete a file"}
	cmdHelp["del"] = `Deletes the file at <path>.

If <rev> is not greater than or equal to the revision of the file,
no change will be made.
`
}


func del(path, rev string) {
	oldRev := mustAtoi64(rev)

	c, err := doozer.Dial(*addr)
	if err != nil {
		bail(err)
	}

	err = c.Del(path, oldRev)
	if err != nil {
		bail(err)
	}
}
