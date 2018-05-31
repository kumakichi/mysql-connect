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
	sshPort    int
	sshHost    string
	sshUser    string
	mysqlPath  string
	sshPath    string
)

func init() {
	parseArgs()

	groups = make(map[string]map[string]string)
	initConfigurations()
	readGroups()
}

func parseArgs() {
	if len(os.Args) < 2 {
		usage()
	}

	flag.IntVar(&sshPort, "p", 22, "SSH port")
	flag.StringVar(&sshHost, "host", "", "SSH host name/ip")
	flag.StringVar(&sshUser, "u", "root", "SSH user name")
	flag.StringVar(&sshPath, "s", "ssh", "ssh program path")
	flag.StringVar(&mysqlPath, "m", "mysql", "mysql program path")
	flag.Parse()
}

func main() {
	fags := flag.Args()
	cmd := fags[0]
	args := fags[1:]

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
