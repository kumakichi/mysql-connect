package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
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

func rename(args []string) {
	if len(args) != 2 {
		errArgs("Should be like: %s mv groupNameOld groupNameNew\n", os.Args[0])
	}

	if grp, ok := groups[args[0]]; ok {
		groups[args[1]] = grp
		delete(groups, args[0])
		updateMyCnf()
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

func dump(args []string) {
	if len(args) < 1 {
		errArgs("Should be like: %s dump group [OPTIONS] [tables]\n", os.Args[0])
	}

	if _, ok := groups[args[0]]; !ok {
		log.Fatalf("Group %s not exists\n", args[0])
	} else {
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
			_, opts := genMysqlCmd("", groups[args[0]])
			dumpOpts := append(opts, args[1:]...)
			exec_command(mysqlDumpPath, dumpOpts...)
		} else {
			cmd, _ := genMysqlCmd(mysqlDumpPath, groups[args[0]])
			dumpCmd := fmt.Sprintf("%s %s", cmd, strings.Join(args[1:], " "))
			exec_command(sshPath, fmt.Sprintf("-p %s", sshPort), "-t", fmt.Sprintf("%s@%s", sshUser, sshHost), dumpCmd)
		}
	}
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
		cmd, _ := genMysqlCmd(mysqlPath, groups[args[0]])
		if identity_file, ok := group["ssh_identity_file"]; ok {
			exec_command(sshPath, "-i", identity_file, fmt.Sprintf("-p %s", sshPort), "-t", fmt.Sprintf("%s@%s", sshUser, sshHost), cmd)
		} else {
			exec_command(sshPath, fmt.Sprintf("-p %s", sshPort), "-t", fmt.Sprintf("%s@%s", sshUser, sshHost), cmd)
		}
	}
}

func genMysqlCmd(program string, group map[string]string) (string, []string) {
	opts := []string{}

	if user, ok := group["user"]; ok {
		userOpt := fmt.Sprintf("-u%s", user)
		opts = append(opts, userOpt)
	}

	if password, ok := group["password"]; ok {
		passOpt := fmt.Sprintf("-p%s", password)
		opts = append(opts, passOpt)
	}

	if host, ok := group["host"]; ok {
		hostOpt := fmt.Sprintf("-h%s", host)
		opts = append(opts, hostOpt)
	}

	if database, ok := group["database"]; ok {
		dbOpt := fmt.Sprintf("%s", database)
		opts = append(opts, dbOpt)
	}

	program += " " + strings.Join(opts, " ")
	return program, opts
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
