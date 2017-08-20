package main

import (
	"flag"
	"github.com/coreos/go-systemd/daemon"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"
)

type webPage struct {
	Title  string
	Header string
	Body   []byte
}

// TODO use a channel to report an error
func keepAlive() {
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil || interval == 0 {
		return
	}
	for {
		_, err := http.Get("http://127.0.0.1:8001")
		if err == nil {
			_, err = daemon.SdNotify(false, "WATCHDOG=1")
			log.Println("notified")
			if err != nil {
				log.Panicf("unable to notify systemd %s", err)
			}
		}
		time.Sleep(interval / 3)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	page := webPage{
		"testing webpage",
		"test header",
		[]byte("some bytes to test body"),
	}

	temp, err := template.ParseFiles("templates/index.tmpl")
	if err != nil {
		log.Panicln("unable to parse template")
	}

	err = temp.Execute(w, page)
	if err != nil {
		log.Panicln("unable to execute template")
	}
}

func main() {
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panicf("cannot listen: %s", err)
	}

	// notify systemd
	_, err = daemon.SdNotify(false, "READY=1")
	if err != nil {
		log.Panicln("unable to notify systemd")
	}

	// keep systemd service alive via watchdog
	go keepAlive()

	var dir string

	flag.StringVar(&dir, "dir", "./static", "directory for serving files")

	router := mux.NewRouter()

	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8081",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	router.HandleFunc("/", handleRoot)

	err = srv.Serve(l)
	if err != nil {
		log.Panicln(err)
	}
}
