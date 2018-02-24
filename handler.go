package main

import (
	"net/http"
	log "github.com/sirupsen/logrus"
	"strings"
	"io/ioutil"
)

type Handler struct {
	FileLogPath string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	reqlog := log.WithFields(
		log.Fields{
			"ip": ip,
			"method": r.Method,
			"url": r.URL.String(),
			"agent": r.UserAgent(),
		})
	if r.Method != "GET" {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		reqlog.Info(http.StatusText(code), code)
		return
	}

	if r.URL.Path == "/log" {
		http.ServeFile(w, r, h.FileLogPath)
		reqlog.Info("return", h.FileLogPath)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/rpc") {
		resp, err := http.Get("http://localhost%v" + r.URL.String())
		var buffer []byte
		if err == nil {
			buffer, err = ioutil.ReadAll(resp.Body)
		}
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			reqlog.WithError(err).Error(http.StatusText(code), code)
			return
		}

		w.Write(buffer)
		reqlog.Info("OK")
		return
	}

	code := http.StatusNotFound
	http.Error(w, http.StatusText(code), code)
	reqlog.Error(http.StatusText(code), code)
}
