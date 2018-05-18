package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
	}

	switch os.Args[1] {
	case "get":
		get()
	case "list":
		list()
	case "del":
		delete()
	case "set":
		set()
	case "conn":
		connect()
	default:
		fmt.Printf("Unsupported command: %s\n", os.Args[1])
		os.Exit(-1)
	}
}

func usage() {
	fmt.Printf("Usage: %s [get/set/list/delete/connect] [group]\n", os.Args[0])
	os.Exit(0)
}
