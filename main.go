package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	myCnf      string
	groups     map[string]map[string]string
	versionStr string
	mysqlPath  string
	sshPath    string
	help       bool
	version    bool
)

func init() {
	parseArgs()

	groups = make(map[string]map[string]string)
	initConfigurations()
	readGroups()
}

func parseArgs() {
	flag.StringVar(&sshPath, "s", "ssh", "ssh program path")
	flag.StringVar(&mysqlPath, "m", "mysql", "mysql program path")
	flag.BoolVar(&help, "h", false, "Print this help infomation")
	flag.BoolVar(&version, "v", false, "Print version")
	flag.Parse()

	if help {
		showVersion()
		showUsageAndExit()
	}

	if version {
		showVersion()
		os.Exit(0)
	}
}

func main() {
	arguments := flag.Args()
	if len(arguments) < 1 {
		showUsageAndExit()
	}

	cmd := arguments[0]
	args := arguments[1:]

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
	default:
		fmt.Printf("Unsupported command: %s\n", os.Args[1])
		os.Exit(-1)
	}
}

func showUsageAndExit() {
	fmt.Printf("Usage: %s [options] command\n", os.Args[0])
	fmt.Printf("command list: [set/ls/show/del/conn/add/cp/fmt/delo]\n")
	fmt.Printf("options: \n")
	flag.PrintDefaults()
	os.Exit(0)
}

func showVersion() {
	fmt.Printf("%s, version: %s\n", os.Args[0], versionStr)
}
