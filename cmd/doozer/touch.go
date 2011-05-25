package main


func init() {
	cmds["touch"] = cmd{touch, "<path>", "update rev of a file"}
	cmdHelp["touch"] = `Attempts to update the rev of a file to a value greater than
the current rev. If a file does not exist, it will be created.
`
}


func touch(path string) {
	c := dial()
	body, rev, err := c.Get(path, nil)
	if err != nil {
		bail(err)
	}
	_, err = c.Set(path, rev, body)
	if err != nil {
		bail(err)
	}
}
