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
	Title string
	Body  []byte
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
			if err != nil {
				log.Panicf("unable to notify systemd %s", err)
			}
		}
		time.Sleep(interval / 3)
	}
}

func handleBaseTemplate(
	tmpl *template.Template,
	data interface{},
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "base", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func main() {
	var err error
	templates := make(map[string]*template.Template)

	templates["index"], err = template.New("index").ParseFiles(
		"./templates/base.tmpl",
		"./templates/index.tmpl",
		"./templates/nav.tmpl",
	)
	if err != nil {
		log.Panicf("unable to load index template: %s", err)
		return
	}

	templates["wip"], err = template.New("wip").ParseFiles(
		"./templates/base.tmpl",
		"./templates/nav.tmpl",
		"./templates/wip.tmpl",
	)
	if err != nil {
		log.Panicf("unable to load wip template: %s", err)
		return
	}

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

	router.HandleFunc("/", handleBaseTemplate(templates["index"], nil))
	router.HandleFunc("/music", handleBaseTemplate(templates["wip"], nil))
	router.HandleFunc("/tech", handleBaseTemplate(templates["wip"], nil))
	router.HandleFunc("/store", handleBaseTemplate(templates["wip"], nil))
	router.HandleFunc("/contact", handleBaseTemplate(templates["wip"], nil))
	router.HandleFunc("/about", handleBaseTemplate(templates["wip"], nil))

	err = srv.Serve(l)
	if err != nil {
		log.Panicln(err)
	}
}
