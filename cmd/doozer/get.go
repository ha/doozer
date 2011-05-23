package main

import (
	"os"
	"fmt"
)


func init() {
	cmds["get"] = cmd{get, "<path>", "read a file"}
	cmdHelp["get"] = "Prints the body and revision of the file at <path>.\n"
}


func get(path string) {
	c := dial()

	body, actualRev, err := c.Get(path, nil)
	if err != nil {
		bail(err)
	}

	fmt.Println("Revision:", actualRev)
	os.Stdout.Write(body)
}
