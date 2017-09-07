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

var defaultStaticDir = "./static/"

type webPage struct {
	Title string
	Body  []byte
}

// KeepAlive uses systemd watchdog to keep
// the server alive
// TODO use a channel to report an error
func KeepAlive() {
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

// ParseFlags parses command line flags
func ParseFlags(staticDir *string) {
	flag.StringVar(staticDir, "staticDir", defaultStaticDir,
		"directory for serving static files")

	flag.Parse()
}

func main() {
	var err error
	templates := make(map[string]*template.Template)

	var staticDir string
	ParseFlags(&staticDir)

	templates["index"], err = template.New("index").ParseFiles(
		staticDir+"templates/base.tmpl",
		staticDir+"templates/index.tmpl",
		staticDir+"templates/nav.tmpl",
	)
	if err != nil {
		log.Panicf("unable to load index template: %s", err)
		return
	}

	templates["wip"], err = template.New("wip").ParseFiles(
		staticDir+"templates/base.tmpl",
		staticDir+"templates/nav.tmpl",
		staticDir+"templates/wip.tmpl",
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
	go KeepAlive()

	router := mux.NewRouter()

	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// TODO change port to something not in go docs
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
