package main

import (
	// "encoding/json"
	"flag"
	"github.com/coreos/go-systemd/daemon"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var defaultStaticDir = "./static/"

type webPage struct {
	Title string
	Body  []byte
}

type contactInfo struct {
	Name  string
	Phone string
	Email string
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
		err := tmpl.ExecuteTemplate(w, "base", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type emailForm struct {
	Name    string
	Email   string
	Phone   string
	Message string
}

func handleEmailSend(email string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// decoder := json.NewDecoder(r.Body)

		// var form emailForm
		// err := decoder.Decode(&form)
		// if err != nil {
		// 	log.Panicf("could not decode message:\n%s", r.Body)
		// }

		// log.Printf("struct formed from json:\n%s\n", form)

		// if email == "none" {
		// 	log.Println("email address is not valid or was not given\n" +
		// 		"unable to send email")
		// } else {

		// }

		log.Printf("form:%v\n", r.PostForm)
	}
}

// ParseFlags parses command line flags
// returns the directory where static files are served
// and the admin email respectively
// TODO add flags for email address and password
func ParseFlags() (string, []byte) {
	var staticDir string
	flag.StringVar(&staticDir, "staticDir", defaultStaticDir,
		"directory for serving static files")

	// flag.StringVar(&key, "key", "", "csrf key")

	// if len(key) != 32 {
	// 	log.Panicf("key is not valid:\n%s\n", key)
	// }

	flag.Parse()

	return staticDir, nil
}

func loadRegularTemplate(name string,
	templateDir string) *template.Template {

	path := templateDir + name + ".tmpl"

	tmpl, err := template.New(name).ParseFiles(
		templateDir+"base.tmpl",
		templateDir+"nav.tmpl",
		templateDir+name+".tmpl",
	)
	if err != nil {
		log.Panicf("could not load %s template\npath: %s\nerr: %s",
			name, path, err)
	}
	return tmpl
}

func sendMail(form emailForm) error {
	return nil
}

func handleSignals() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGABRT)
	sig := <-sigChannel
	log.Printf("signal recieved: %s\n", sig)
	panic("process aborted")
}

func main() {
	var err error
	templates := make(map[string]*template.Template)

	// TODO encode contact info
	//////////////////////////////////////////////
	// myInfo := contactInfo{				    //
	// 	Name:  "Chris Overcash",			    //
	// 	Email: "covercash.biz@gmail.com",	    //
	// 	Phone: "(501) 510-0946",			    //
	// }									    //
	//////////////////////////////////////////////

	staticDir, key := ParseFlags()

	templateDir := staticDir + "templates/"

	templateNames := [...]string{
		"about",
		"contact",
		"index",
		"music",
		"tech",
		"wip",
	}

	for _, s := range templateNames {
		templates[s] = loadRegularTemplate(s, templateDir)
	}

	listener, err := net.Listen("tcp", ":8081")
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
	go handleSignals()

	router := mux.NewRouter()

	router.PathPrefix("/static/").
		Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	key = []byte("000000000TEST0000000000000000000")

	handler := csrf.Protect(key, csrf.Secure(false))(router)

	// TODO change port to something not in go docs
	srv := &http.Server{
		Handler:      handler,
		Addr:         "127.0.0.1:8081",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	apiRouter := router.PathPrefix("/api").Subrouter()

	router.HandleFunc("/", handleBaseTemplate(templates["index"], nil))
	router.HandleFunc("/music", handleBaseTemplate(templates["music"], nil))
	router.HandleFunc("/tech", handleBaseTemplate(templates["tech"], nil))
	router.HandleFunc("/store", handleBaseTemplate(templates["wip"], nil))
	router.HandleFunc("/contact", handleBaseTemplate(templates["contact"], nil))
	router.HandleFunc("/about", handleBaseTemplate(templates["about"], nil))

	// TODO fix contact info
	email := "chris@covercash.biz"
	// router.HandleFunc("/api/email", handleEmailSend(email))
	apiRouter.HandleFunc("/email", handleEmailSend(email)).Methods("POST")

	err = srv.Serve(listener)
	if err != nil {
		log.Panicln(err)
	}
}
