package main

import (
	"flag"
	"net/http"
	"fmt"
)

const (
	LICENSER_PORT = 8080
)

var (
	port = flag.Int("p", 80, "port")
	user = flag.String("u", "wolfogre.com", "user")
	redirect = flag.String("r", "http://blog.wolfogre.com/posts/jetbrains-licenser/", "redirect")
	logpath = flag.String("l", "/opt/log/last.log", "log path")
	binpath = flag.String("b", "/opt/bin/jetbrains-licenser", "bin path")
)

func main() {
	InitLog()
	defer SyncLog()

	flag.Parse()
	InitFileLog(*logpath)
	Log.Infow("start",
			"port", *port,
			"user", *user,
			"redirect", *redirect,
			"logpath", *logpath,
			"binpath", *binpath)
	go Licenser(*binpath, LICENSER_PORT, *user)

	http.ListenAndServe(fmt.Sprintf(":%v", *port), &Handler{
		FileLogPath: *logpath,
		RedirectUrl: *redirect,
		LicenserAddr: fmt.Sprintf("http://localhost:%v", LICENSER_PORT),
	})
}
