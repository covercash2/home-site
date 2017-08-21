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

var templates = map[string]*template.Template{
	"index": template.Must(template.ParseFiles(
		"templates/base.tmpl",
		"templates/index.tmpl",
		"templates/nav.tmpl",
	)),
	"wip": template.Must(template.ParseFiles(
		"templates/base.tmpl",
		"templates/nav.tmpl",
		"templates/wip.tmpl",
	)),
}

type webPage struct {
	Title string
	Body  []byte
}

var indexPage = webPage{
	"C/$",
	[]byte("test page"),
}

func renderTemplate(w http.ResponseWriter, tmpl string, page *webPage) {
	err := templates[tmpl].ExecuteTemplate(w, "base", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderPlainTemplate(w http.ResponseWriter, tmpl string) {
	err := templates[tmpl].ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	page := webPage{
		"testing webpage",
		[]byte("some bytes to test body"),
	}

	renderTemplate(w, "index", &page)
}

func handleWip(w http.ResponseWriter, r *http.Request) {
	renderPlainTemplate(w, "wip")
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
	router.HandleFunc("/music", handleWip)
	router.HandleFunc("/tech", handleWip)
	router.HandleFunc("/store", handleWip)
	router.HandleFunc("/contact", handleWip)
	router.HandleFunc("/about", handleWip)

	err = srv.Serve(l)
	if err != nil {
		log.Panicln(err)
	}
}
