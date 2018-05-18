package main

import (
	"fmt"
	"log"
	"os"
)

func list() {
	for key, _ := range groups {
		fmt.Println(key)
	}
}

func show(args []string) {
	if len(args) != 1 {
		errArgs("Should be like: %s show groupName1 [groupName2]...\n", os.Args[0])
	}

	for _, val := range args {
		fmt.Printf(" ------- %s -------\n", val)
		if g, ok := groups[val]; ok {
			printGroup(&g)
		} else {
			fmt.Println("not found")
		}
	}
}

func del(args []string) {
	if len(args) < 1 {
		errArgs("Should be like: %s del group1 [group2] ...\n", os.Args[0])
	}

	for _, v := range args {
		delete(groups, v)
		updateMyCnf()
	}
}

func set(args []string) {
	if len(args) < 2 {
		errArgs("Should be like: %s set groupName host=someHost user=someUser ...\n", os.Args[0])
	}

	if group, ok := groups[args[0]]; ok {
		updateBody(group, args[1:])
	} else {
		log.Fatalf("Group %s not found\n", args[0])
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
