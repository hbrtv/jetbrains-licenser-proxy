package main

import (
	"bufio"
	"github.com/bitly/go-simplejson"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	resultM = &sync.RWMutex{}
	dateResult []string
	userSet = make(map[string]bool)
	timesResult = make(map[string]int)
	userResult = make(map[string]map[string]bool)
	newUserResult = make(map[string]int)
	ipResult = make(map[string]map[string]bool)
	productResult = make(map[string]map[string]bool)
)

func InitStatistics(fileLogPath string) {
	resultM.Lock()
	defer resultM.Unlock()

	f, err := os.Open(fileLogPath)
	if err != nil{
		Log.Errorf("failed to open %v", fileLogPath)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		if count >= 10000 && count % 10000 == 0 {
			runtime.GC()
			Log.Infof("smooth starting: %v", count)
			time.Sleep(time.Second)
		}
		count++
		if len(scanner.Bytes()) == 0 {
			continue
		}
		j, err := simplejson.NewJson(scanner.Bytes())
		if err != nil {
			Log.Errorf("failed to new json: %v, %v", err, scanner.Text())
			return
		}
		date := j.Get("time").MustString()
		if date == "" {
			continue
		}
		date = strings.Replace(date[:10], "-", "", -1)
		appendDate(date)

		times, _ := timesResult[date]
		timesResult[date] = times + 1

		user, ok := userResult[date]
		if !ok {
			user = make(map[string]bool)
		}
		user[j.Get("machineId").MustString()] = true
		userResult[date] = user

		if _, ok := userSet[j.Get("machineId").MustString()]; !ok {
			userSet[j.Get("machineId").MustString()] = true
			newUser, _ := newUserResult[date]
			newUserResult[date] = newUser + 1
		}

		ip, ok := ipResult[date]
		if !ok {
			ip = make(map[string]bool)
		}
		ip[j.Get("ip").MustString()] = true
		ipResult[date] = ip

		product, ok := productResult[date]
		if !ok {
			product = make(map[string]bool)
		}
		product[j.Get("productCode").MustString()] = true
		productResult[date] = product
	}
	runtime.GC()
	Log.Infof("smooth started: %v", count)
}

func GetStatistics() map[string]interface{} {
	resultM.RLock()
	defer resultM.RUnlock()

	var times, user, newUser, ip, product []int
	for _, date := range dateResult{
		if v, ok := timesResult[date]; ok {
			times = append(times, v)
		} else {
			times = append(times, 0)
		}
		if v, ok := userResult[date]; ok {
			user = append(user, len(v))
		} else {
			user = append(user, 0)
		}
		if v, ok := newUserResult[date]; ok {
			newUser = append(newUser, v)
		} else {
			newUser = append(newUser, 0)
		}
		if v, ok := ipResult[date]; ok {
			ip = append(ip, len(v))
		} else {
			ip = append(ip, 0)
		}
		if v, ok := productResult[date]; ok {
			product = append(product, len(v))
		} else {
			product = append(product, 0)
		}
	}
	return map[string]interface{}{
		"Date": dateResult,
		"Times": times,
		"User": user,
		"NewUser": newUser,
		"IP": ip,
		"Product": product,
	}
}

func AppendStatistics(time, machineId, _ip, productCode string) {
	resultM.Lock()
	defer resultM.Unlock()

	date := strings.Replace(time[:10], "-", "", -1)
	appendDate(date)

	times, _ := timesResult[date]
	timesResult[date] = times + 1

	user, ok := userResult[date]
	if !ok {
		user = make(map[string]bool)
	}
	user[machineId] = true
	userResult[date] = user

	if _, ok := userSet[machineId]; !ok {
		userSet[machineId] = true
		newUser, _ := newUserResult[date]
		newUserResult[date] = newUser + 1
	}

	ip, ok := ipResult[date]
	if !ok {
		ip = make(map[string]bool)
	}
	ip[_ip] = true
	ipResult[date] = ip

	product, ok := productResult[date]
	if !ok {
		product = make(map[string]bool)
	}
	product[productCode] = true
	productResult[date] = product
}

func appendDate(start string) {
	if len(dateResult) == 0 {
		dateResult = append(dateResult, start)
	}
	for {
		newest, err := time.ParseInLocation("20060102", dateResult[len(dateResult) - 1], time.Local)
		if err != nil {
			Log.Panic(err)
		}
		now := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		if !now.After(newest) {
			break
		}
		dateResult = append(dateResult, newest.AddDate(0, 0, 1).Format("20060102"))
	}
}
