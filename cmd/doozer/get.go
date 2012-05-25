package main

import (
	"os"
	"fmt"
)


func init() {
	cmds["get"] = cmd{get, "<path> <rev|nil>", "read a file"}
	cmdHelp["get"] = "Prints the body and revision of the file at <path>.\n"
}


func get(path string, rev string) {
	c := dial()

	var (
		body      []byte
		actualRev int64
		err       os.Error
	)

	if rev == "nil" {
		body, actualRev, err = c.Get(path, nil)
	} else {
		requestedRev := mustAtoi64(rev)
		body, actualRev, err = c.Get(path, &requestedRev)
	}

	if err != nil {
		bail(err)
	}

	fmt.Println("Revision:", actualRev)
	os.Stdout.Write(body)
}
