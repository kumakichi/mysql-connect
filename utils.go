package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	defaultGroup  = "default"
	clientPrefix  = "client"
	mysqlPrefix   = "mysql"
	userDefPrefix = "usrDef"
)

var (
	regHead  *regexp.Regexp
	regBody  *regexp.Regexp
	curGroup string
)

func initConfigurations() {
	switch runtime.GOOS {
	case "linux":
		myCnf = filepath.Join(os.Getenv("HOME"), ".my.cnf")
		regHead = regexp.MustCompile(`(?U)(\s*)(\[)(.*)(\])(\s*$)`)
		regBody = regexp.MustCompile(`(?U)(\s*)([^\s]+)(=)([^\s]+)(\s*$)`)
	default:
		fmt.Println("Currently, support only linux platform")
		os.Exit(0)
	}
}

func readGroups() {
	file, err := os.Open(myCnf)
	if err != nil {
		fmt.Printf("Read conf file : %s failed\n", myCnf)
		return
	}
	defer file.Close()

	br := bufio.NewReader(file)
	for {
		line, isPrefix, err := br.ReadLine()
		if err != nil {
			if err != io.EOF {
				err = err
			}
			break
		}

		if isPrefix {
			fmt.Println("A too long line, seems unexpected")
			return
		}

		str := strings.TrimSpace(string(line))
		parseLine(str)
	}
}

func parseLine(str string) {
	if regHead.MatchString(str) {
		parseHead(str)
	} else if regBody.MatchString(str) {
		parseBody(str)
	}
}

func parseHead(str string) {
	matches := regHead.FindStringSubmatch(str)
	head := matches[3]

	var typ string
	if strings.HasPrefix(head, clientPrefix) {
		typ = clientPrefix
		curGroup = strings.Replace(head, clientPrefix, "", 1)
	} else if strings.HasPrefix(head, mysqlPrefix) {
		typ = mysqlPrefix
		curGroup = strings.Replace(head, mysqlPrefix, "", 1)
	} else {
		typ = userDefPrefix
		curGroup = head
	}

	if curGroup == "" {
		curGroup = defaultGroup
	}

	if _, ok := groups[curGroup]; !ok {
		groups[curGroup] = make(map[string]string)
	}
	groups[curGroup]["typ"] = typ
}

func parseBody(str string) {
	matches := regBody.FindStringSubmatch(str)
	key := matches[2]
	val := matches[4]

	if key != "user" && key != "password" && key != "host" && key != "database" {
		log.Fatalf("Unsupported key: %s in line [%s]\n", key, str)
	}
	groups[curGroup][key] = val
}

func printGroup(group *map[string]string) {
	for k, v := range *group {
		if v != "" {
			fmt.Printf("%s\t%s\n", k, v)
		}
	}
}
