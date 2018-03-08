package main

import (
	"os"
	"sync"
	"net/url"
	"fmt"
	"time"
	"strings"
	"encoding/json"
	"io/ioutil"
	"bufio"
	"bytes"
)

var (
	file *os.File
	fileM *sync.Mutex
	lastLog [][]byte
	lastLogM *sync.RWMutex
)

func InitFileLog(logpath string) error {
	fileM = &sync.Mutex{}
	lastLogM = &sync.RWMutex{}

	if buffer, err := ioutil.ReadFile(logpath); err == nil{
		scanner := bufio.NewScanner(bytes.NewReader(buffer))
		for scanner.Scan() {
			if len(scanner.Bytes()) == 0 {
				continue
			}
			line := make([]byte, len(scanner.Bytes()))
			copy(line, scanner.Bytes())
			AppendLastLog(line)
		}
	}

	f, err := os.OpenFile(logpath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	file = f
	return nil
}

func FileLog(ip, location string, u *url.URL) {
	m := make(map[string]string)
	t := time.Now().Format("2006-01-02T15:04:05")
	m["time"] = t
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
	AppendLastLog(buffer)
	AppendStatistics(t, u.Query().Get("machineId"), ip, u.Query().Get("productCode"))
}

func AppendLastLog(log []byte) {
	lastLogM.Lock()
	defer lastLogM.Unlock()
	lastLog = append(lastLog, log)
	for len(lastLog) > 100 {
		lastLog = lastLog[1:]
	}
}

func GetLastLog() []byte {
	lastLogM.RLock()
	defer lastLogM.RUnlock()
	result := make([]byte, 0)
	for _, v := range lastLog {
		result = append(result, v...)
		result = append(result, '\n')
	}
	return result
}