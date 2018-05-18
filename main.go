package main

import (
	"fmt"
	"os"
)

var (
	myCnf  string
	groups map[string]map[string]string
)

func main() {
	groups = make(map[string]map[string]string)
	initConfigurations()
	readGroups()

	if len(os.Args) < 2 {
		usage()
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "get":
		get(args)
	case "list":
		list()
	case "show":
		show(args)
	case "del":
		del(args)
	case "set":
		set(args)
	case "conn":
		conn(args)
	default:
		fmt.Printf("Unsupported command: %s\n", os.Args[1])
		os.Exit(-1)
	}
}

func usage() {
	fmt.Printf("Usage: %s [get/set/list/show/del/connect] [group]\n", os.Args[0])
	os.Exit(0)
}
