package main

import (
	"fmt"
	"os"
)

var (
	myCnf      string
	groups     map[string]map[string]string
	versionStr string
)

func init() {
	if len(os.Args) < 2 {
		usage()
	}

	groups = make(map[string]map[string]string)
	initConfigurations()
	readGroups()
}

func main() {
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "ls":
		list()
	case "show":
		show(args)
	case "del":
		del(args)
	case "add":
		add(args)
	case "set":
		set(args)
	case "conn":
		conn(args)
	case "cp":
		cp(args)
	case "fmt":
		format()
	case "delo":
		delOption(args)
	case "-h", "--help":
		usage()
	case "-v", "--version":
		version()
	default:
		fmt.Printf("Unsupported command: %s\n", os.Args[1])
		os.Exit(-1)
	}
}

func usage() {
	fmt.Printf("Usage: %s [set/ls/show/del/conn/add/cp/fmt/delo] [options]\n", os.Args[0])
	os.Exit(0)
}

func version() {
	fmt.Println(versionStr)
}
