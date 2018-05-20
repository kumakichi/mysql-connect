package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const (
	tmpCnfFile    = ".my.cnf.mc"
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
		parseBodyLine(str)
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

func parseBodyLine(str string, onlyValidate ...bool) (string, string) {
	matches := regBody.FindStringSubmatch(str)
	key := matches[2]
	val := matches[4]

	if key != "user" && key != "password" && key != "host" && key != "database" {
		log.Fatalf("Unsupported key: %s in line [%s]\n", key, str)
	}

	if len(onlyValidate) == 0 {
		groups[curGroup][key] = val
		return "", ""
	}

	return key, val
}

func printGroup(group *map[string]string) {
	for k, v := range *group {
		if v != "" {
			fmt.Printf("%-8s\t%s\n", k, v)
		}
	}
}

func updateBody(group map[string]string, args []string) {
	m := make(map[string]string)
	// validate
	for _, line := range args {
		if !regBody.MatchString(line) {
			log.Fatalf("Invalid line: %s\n", line)
		} else {
			k, v := parseBodyLine(line, true)
			m[k] = v
		}
	}
	// do the work
	for key, val := range m {
		group[key] = val
	}

	updateMyCnf()
}

func updateMyCnf() {
	var path string
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := ioutil.TempFile(wd, tmpCnfFile)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range groups {
		writeHead(file, k, v)
		writeBody(file, v)
		file.WriteString("\n")
	}

	path = file.Name()
	file.Close()

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("mv %s %s", path, myCnf))
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func writeHead(file *os.File, str string, m map[string]string) {
	if str == defaultGroup {
		file.WriteString(fmt.Sprintf("[%s]\n", m["typ"]))
	} else {
		file.WriteString(fmt.Sprintf("[%s%s]\n", m["typ"], str))
	}
}

func writeBody(file *os.File, m map[string]string) {
	for k, v := range m {
		if k != "typ" && len(k) > 0 && len(v) > 0 {
			file.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
	}
}
