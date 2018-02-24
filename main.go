package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
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
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.WithFields(
		log.Fields{
			"port": *port,
			"user": *user,
			"redirect": *redirect,
			"logpath": *logpath,
			"binpath": *binpath,
		}).Info("start")
	go Licenser(*binpath, LICENSER_PORT, *user)

	//http.ListenAndServe(fmt.Sprintf(":%v", *port), )
}
