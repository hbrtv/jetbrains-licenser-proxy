package main

import (
	"os"
	"sync"
	"net/url"
	"fmt"
	"time"
	"strings"
	"encoding/json"
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

func FileLog(ip, location string, u *url.URL) {
	m := make(map[string]string)
	m["time"] = time.Now().Format("2006-01-02T15:04:05")
	m["ip"] = ip
	m["location"] = location
	for k, v := range u.Query() {
		m[k] = strings.Join(v, ",")
	}
	buffer, _ := json.Marshal(m)
	fileM.Lock()
	_, err := file.WriteString(fmt.Sprintf("%s\n", buffer))
	fileM.Unlock()
	if err != nil {
		Log.Panic(err)
	}
}
