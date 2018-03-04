package main

import (
	"io/ioutil"
	"bufio"
	"bytes"

	"github.com/bitly/go-simplejson"
	"sync"
	"strings"
	"time"
	"fmt"
)

var (
	resultM = &sync.RWMutex{}
	dateResult []string
	timesResult = make(map[string]int)
	userResult = make(map[string]map[string]bool)
	ipResult = make(map[string]map[string]bool)
	productResult = make(map[string]map[string]bool)
)

func InitStatistics(fileLogPath string) {
	resultM.Lock()
	defer resultM.Unlock()

	buffer, err := ioutil.ReadFile(fileLogPath)
	if err != nil{
		Log.Errorf("failed to read %v", fileLogPath)
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(buffer))
	for scanner.Scan() {
		if len(scanner.Bytes()) == 0 {
			continue
		}
		j, err := simplejson.NewJson(scanner.Bytes())
		if err != nil {
			Log.Errorf("failed to new json: %v, %v", err, scanner.Text())
			return
		}
		date := j.MustString("time")
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
		user[j.MustString("machineId")] = true
		userResult[date] = user

		ip, ok := ipResult[date]
		if !ok {
			ip = make(map[string]bool)
		}
		ip[j.MustString("ip")] = true
		ipResult[date] = ip

		product, ok := productResult[date]
		if !ok {
			product = make(map[string]bool)
		}
		product[j.MustString("productCode")] = true
		productResult[date] = product
	}

}

func GetStatistics() map[string]interface{} {
	var times, user, ip, product []int
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
		"IP": ip,
		"Product": product,
	}
}

func AppendLog(time, machineId, _ip, productCode string) {
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
