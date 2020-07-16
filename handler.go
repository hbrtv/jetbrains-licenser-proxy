package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	FileLogPath  string
	TemplateDir  string
	RedirectUrl  string
	LicenserAddr string
	client       *http.Client
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	reqlog := Log.With(
		"ip", ip,
		"method", r.Method,
		"url", r.URL.String(),
		"agent", r.UserAgent())
	if r.Method != "GET" {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		reqlog.Infof("%v %v", http.StatusText(code), code)
		return
	}

	if r.URL.Path == "/" {
		http.Redirect(w, r, h.RedirectUrl, http.StatusTemporaryRedirect)
		reqlog.Infof("redirect %v", h.RedirectUrl)
		return
	}

	if r.URL.Path == "/log" {
		lastLog := GetLastLog()
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.Header().Add("Content-Length", strconv.Itoa(len(lastLog)))
		w.Write(lastLog)
		reqlog.Infof("log")
		return
	}

	if r.URL.Path == "/statistics" {
		temp, err := template.ParseFiles(h.TemplateDir + "/statistics.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = temp.Execute(w, GetStatistics())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if r.URL.Path == "/_status" {
		w.Write([]byte("OK"))
		return
	}

	if strings.HasPrefix(r.URL.Path, "/rpc") {
		if h.client == nil {
			h.client = &http.Client{
				Timeout: 5 * time.Second,
			}
		}
		resp, err := h.client.Get(h.LicenserAddr + r.URL.String())
		var buffer []byte
		if err == nil {
			buffer, err = ioutil.ReadAll(resp.Body)
		}
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			reqlog.Infof("%v %v", http.StatusText(code), code, err)
			return
		}

		w.Write(buffer)
		reqlog.Info("OK")
		FileLog(ip, h.LocateIP(ip), r.URL)
		return
	}

	code := http.StatusNotFound
	http.Error(w, http.StatusText(code), code)
	reqlog.Infof("%v %v", http.StatusText(code), code)
}

type LocateResponse struct {
	Code int `json:"code"`
	Data struct {
		Country string `json:"country"`
		Area    string `json:"area"`
		Region  string `json:"region"`
		City    string `json:"city"`
		County  string `json:"county"`
		Isp     string `json:"isp"`
	} `json:"data"`
}

func (h *Handler) LocateIP(ip string) string {
	log := Log.With("ip", ip)
	resp, err := h.client.Get("http://ip.taobao.com/service/getIpInfo.php?accessKey=alibaba-inc&ip=" + ip)
	if err != nil {
		log.Infof("can not get ip location: %v", err)
		return ""
	}
	buffer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Infof("can not get ip location: %v", err)
		return ""
	}
	response := LocateResponse{}
	err = json.Unmarshal(buffer, &response)
	if err != nil {
		log.Infof("json unmarshal failed: %v, %s", err, buffer)
		return ""
	}
	if response.Code != 0 {
		log.Infof("can not get ip location: %s", buffer)
		return ""
	}
	data := response.Data
	return fmt.Sprintf("%v,%v,%v,%v,%v,%v", data.Country, data.Area, data.Region, data.City, data.County, data.Isp)
}
