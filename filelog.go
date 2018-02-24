package main

import (
	"os"
	"sync"
	"net/url"
	"fmt"
)

var (
	file *os.File
	fileM *sync.Mutex
)

func InitFileLog(logpath string) error {
	f, err := os.OpenFile(logpath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	file = f
	fileM = &sync.Mutex{}
	return nil
}

func FileLog(ip string, u url.URL) error {
	fileM.Lock()
	defer fileM.Unlock()
	_, err := file.WriteString(fmt.Sprintf("[%v] %v\n", u.String(), u.String()))
	return err
}
