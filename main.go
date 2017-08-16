package main

import (
	"log"
	"net"
	"net/http"
	"time"
	"github.com/coreos/go-systemd/daemon"
)

func keep_alive() error {
	interval, err := daemon.SdWatchdogEnabled(false)
	if err != nil || interval == 0 {
		return err
	}
	for {
		_, err := http.Get("http://127.0.0.1:8001")
		if err == nil {
			daemon.SdNotify(false, "WATCHDOG=1")
		}
		time.Sleep(interval / 3)
	}
	return nil
}

func main() {
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panicf("cannot listen: %s", err)
	}
	daemon.SdNotify(false, "READY=1")

	// keep alive via watchdog
	go keep_alive()

	http.Serve(l, nil)
}
