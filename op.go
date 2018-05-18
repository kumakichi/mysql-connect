package main

import (
	"fmt"
	"os"
)

func get(args []string) {
	if len(args) != 1 {
		errArgs("Should be like: %s get groupName\n", os.Args[0])
	}
}

func list() {

}

func del(args []string) {
	if len(args) < 1 {
		errArgs("Should be like: %s del group1 [group2] ...\n", os.Args[0])
	}
}

func set(args []string) {
	if len(args) < 2 {
		errArgs("Should be like: %s set groupName host=someHost user=someUser ...\n", os.Args[0])
	}
}

func conn(args []string) {
	if len(args) != 1 {
		errArgs("Should be like: %s conn groupName\n", os.Args[0])
	}

	fmt.Printf("mysql --defaults-group-suffix=%s\n", args[0])
}

func errArgs(format string, a ...interface{}) {
	fmt.Println("Invalid arg[s]")
	fmt.Printf(format, a...)
	os.Exit(-1)
}
