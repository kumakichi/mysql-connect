package main

import "os/exec"

func get() {

}

func list() {

}

func delete() {

}

func set() {

}

func connect() {
	exec.Command("mysql", "--defaults-group-suffix")
}
