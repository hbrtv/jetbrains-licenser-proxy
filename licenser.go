package main

import (
	"os/exec"
	"strconv"
	log "github.com/sirupsen/logrus"
	"os"
)

func Licenser(binpath string, port int, user string) {
	log.Info("licenser start")
	cmd := exec.Command(binpath, "-p", strconv.Itoa(port), "-u", user)
	err := cmd.Run()
	log.WithError(err).Error("licenser exit")
	if err == nil {
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
