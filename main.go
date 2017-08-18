package main

import (
	"fmt"
	"github.com/coreos/go-systemd/daemon"
	"log"
	"net"
	"net/http"
	"time"
)

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
	fmt.Fprintf(w, "hi there, I love %s!", r.URL.Path[1:])
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

	http.HandleFunc("/", handleRoot)

	err = http.Serve(l, nil)
	if err != nil {
		log.Panicln(err)
	}
}
