package main

import "fmt"

func init() {
	cmds["help"] = cmd{help, "<command>", "provide detailed help on a command"}
}

func help(command string) {
	c := cmds[command]
	fmt.Printf("Use: %s [options] %s [options] %s\n\n", selfName, command, c.a)
	fmt.Print(cmdHelp[command])
}
