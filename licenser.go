package main

import (
	"os/exec"
	"strconv"
	"os"
)

func Licenser(binpath string, port int, user string) {
	Log.Info("licenser start")
	cmd := exec.Command(binpath, "-p", strconv.Itoa(port), "-u", user)
	err := cmd.Run()
	Log.Error("licenser exit", err)
	if err != nil {
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
