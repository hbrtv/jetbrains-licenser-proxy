package main

import (
	"bufio"
	"container/list"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	file *os.File
	fileM *sync.Mutex
	lastLog *list.List
	lastLogM *sync.RWMutex
)

func InitFileLog(logpath string) error {
	fileM = &sync.Mutex{}
	lastLog = list.New()
	lastLogM = &sync.RWMutex{}

	f, err := os.Open(logpath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		AppendLastLog(scanner.Text())
	}
	f.Close()

	f, err = os.OpenFile(logpath, os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0666)
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
	AppendLastLog(string(buffer))
	AppendStatistics(t, u.Query().Get("machineId"), ip, u.Query().Get("productCode"))
}

func AppendLastLog(log string) {
	lastLogM.Lock()
	defer lastLogM.Unlock()
	lastLog.PushBack(log)
	for lastLog.Len() > 100 {
		lastLog.Remove(lastLog.Front())
	}
}

func GetLastLog() []byte {
	lastLogM.RLock()
	defer lastLogM.RUnlock()
	result := strings.Builder{}
	it := lastLog.Front()
	for it != nil {
		result.WriteString(it.Value.(string))
		result.WriteString("\n")
		it = it.Next()
	}
	return []byte(result.String())
}