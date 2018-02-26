package main

import (
	"net/http"
	"strings"
	"io/ioutil"
	"time"
)

type Handler struct {
	FileLogPath string
	RedirectUrl string
	LicenserAddr string
	client *http.Client
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
		reqlog.Info(http.StatusText(code), code)
		return
	}

	if r.URL.Path == "/" {
		http.Redirect(w, r, h.RedirectUrl, http.StatusTemporaryRedirect)
		reqlog.Infof("redirect %v", h.RedirectUrl)
		return
	}

	if r.URL.Path == "/log" {
		http.ServeFile(w, r, h.FileLogPath)
		reqlog.Info("return", h.FileLogPath)
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
		FileLog(ip, r.URL)
		return
	}

	code := http.StatusNotFound
	http.Error(w, http.StatusText(code), code)
	reqlog.Infof("%v %v", http.StatusText(code), code)
}
