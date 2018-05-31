package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
		fmt.Printf("-------- %s --------\n", val)
		if g, ok := groups[val]; ok {
			printGroup(&g)
		} else {
			fmt.Println("not found")
		}
	}
}

func format() {
	updateMyCnf()
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

func cp(args []string) {
	if len(args) != 2 {
		errArgs("Should be like: %s cp fromGroup toGroup\n", os.Args[0])
	}

	if _, ok := groups[args[0]]; !ok {
		log.Fatalf("Group %s not exists\n", args[0])
	} else {
		groups[args[1]] = make(map[string]string)
		for k, v := range groups[args[0]] {
			groups[args[1]][k] = v
		}
		updateMyCnf()
	}
}

func add(args []string) {
	if len(args) == 0 {
		errArgs("Should be like: %s add groupName [host=someHost] [user=someUser] ...\n", os.Args[0])
	}

	if _, ok := groups[args[0]]; ok {
		log.Fatalf("Group %s already exists\n", args[0])
	} else {
		groups[args[0]] = make(map[string]string)
		groups[args[0]]["typ"] = clientPrefix
		updateBody(groups[args[0]], args[1:])
	}
}

func set(args []string) {
	if len(args) < 2 {
		errArgs("Should be like: %s set groupName host=someHost [user=someUser] ...\n", os.Args[0])
	}

	if group, ok := groups[args[0]]; ok {
		updateBody(group, args[1:])
	} else {
		log.Fatalf("Group %s not found\n", args[0])
	}
}

func delOption(args []string) {
	if len(args) < 2 {
		errArgs("Should be like: %s delo groupName keyName1 [keyName2] ...\n", os.Args[0])
	}

	if group, ok := groups[args[0]]; ok {
		for _, v := range args[1:] {
			if _, ok := group[v]; ok {
				delete(group, v)
			} else {
				fmt.Printf("Key %s not found, fall through...\n", v)
			}
		}
	} else {
		log.Fatalf("Group %s not found\n", args[0])
	}
	updateMyCnf()
}

func conn(args []string) {
	if len(args) != 1 {
		errArgs("Should be like: %s conn groupName\n", os.Args[0])
	}
	groupName := args[0]
	group := groups[groupName]

	sshHost := ""
	sshUser := "root"
	sshPort := "22"

	if host, ok := group["ssh_host"]; ok {
		sshHost = host
	}

	if user, ok := group["ssh_user"]; ok {
		sshUser = user
	}

	if port, ok := group["ssh_port"]; ok {
		sshPort = port
	}

	if sshHost == "" {
		exec_command(mysqlPath, fmt.Sprintf("--defaults-group-suffix=%s", groupName))
	} else {
		cmd := genMysqlConnCmd(groups[args[0]])
		exec_command(sshPath, fmt.Sprintf("-p %s", sshPort), "-t", fmt.Sprintf("%s@%s", sshUser, sshHost), cmd)
	}
}

func genMysqlConnCmd(group map[string]string) string {
	cmd := mysqlPath
	if user, ok := group["user"]; ok {
		cmd += fmt.Sprintf(" -u%s", user)
	}

	if password, ok := group["password"]; ok {
		cmd += fmt.Sprintf(" -p%s", password)
	}

	if host, ok := group["host"]; ok {
		cmd += fmt.Sprintf(" -h%s", host)
	}

	if database, ok := group["database"]; ok {
		cmd += fmt.Sprintf(" %s", database)
	}

	return cmd
}

func errArgs(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(-1)
}

func exec_command(program string, args ...string) {
	cmd := exec.Command(program, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
