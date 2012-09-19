package main

func init() {
	cmds["nop"] = cmd{nop, "", "consensus"}
	cmdHelp["nop"] = `Performs a consensus operation.

No change will be made to the data store.
`
}

func nop() {
	c := dial()

	err := c.Nop()
	if err != nil {
		bail(err)
	}
}
