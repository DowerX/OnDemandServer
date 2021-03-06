package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/DowerX/OnDemandServer/data"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Credentials _
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

var endtime time.Time
var stepsize time.Duration
var limit time.Duration
var looking bool = false
var c data.Config

func main() {
	var conf string
	flag.StringVar(&conf, "config", "./config.yml", "Set the config file.")
	flag.Parse()
	c = data.GetConfig(conf)
	checkFlags()
	endtime = time.Now()
	stepsize = time.Duration(c.Stepsize) * time.Second
	limit = time.Duration(c.Limit) * time.Second

	if c.Log {
		lf, err := os.OpenFile(c.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("ERR: Can't log!")
			return
		}
		log.SetOutput(lf)
	}

	r := mux.NewRouter()
	r.HandleFunc(c.Path, postReq).Methods("POST")
	http.ListenAndServe(c.Port, r)
}

func checkFlags() {
	flag.IntVar(&c.Stepsize, "stepsize", c.Stepsize, "Set stepsize of additions.")
	flag.IntVar(&c.Limit, "limit", c.Limit, "Set the maximum time for the countdown.")
	flag.StringVar(&c.Path, "path", c.Path, "Set the routing path for the server.")
	flag.StringVar(&c.Port, "port", c.Port, "Set the port of the server.")
	flag.StringVar(&c.StartScript, "start", c.StartScript, "Set the starting script of the service.")
	flag.StringVar(&c.StopScript, "stop", c.StopScript, "Set the stopping script of the service.")
	flag.BoolVar(&c.Log, "log", c.Log, "Enable or disbale logging to a file.")
	flag.StringVar(&c.Logfile, "logfile", c.Logfile, "Set the logfile.")
	flag.StringVar(&c.Users, "users", c.Users, "Set the user's credentials file.")
	flag.Parse()
}

func writelog(l ...interface{}) {
	fmt.Println(l, "")
	if c.Log {
		log.Println(l, "")
	}
}

func postReq(w http.ResponseWriter, r *http.Request) {

	creds := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		writelog("WARN: Incoming request, but wrong format!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users := data.GetUsers(c.Users)

	for i := 0; i < len(users); i++ {
		if users[i].Username == creds.Username {
			c := bcrypt.CompareHashAndPassword([]byte(users[i].PasswordHash), []byte(creds.Password))
			if c == nil {
				go req(users[i])
				writelog("INFO: Incoming request:", creds.Username)
				w.WriteHeader(http.StatusAccepted)
				return
			}
		}
	}
	writelog("WARN: Unauthorized request: ", creds.Username)
	w.WriteHeader(http.StatusUnauthorized)
}

func req(user data.User) {
	if endtime.Sub(time.Now()) <= 0 {
		endtime = time.Now()
	}

	if endtime.Sub(time.Now()) <= limit {
		endtime = endtime.Add(stepsize)
	}

	writelog("INFO: Added", stepsize, ".", "Service will stop at", endtime, ".")

	if looking == false {
		//START
		cmd := exec.Command(c.StartScript)
		cmd.Start()
		writelog("INFO: Started service.")
		lookForEnd()
	}
}

func lookForEnd() {
	looking = true
	cmd := exec.Command(c.StopScript)
	for {
		if endtime.Sub(time.Now()) <= 0 {
			//STOP
			cmd.Start()
			writelog("INFO: Stopped service.")
			looking = false
			return
		}
		time.Sleep(10 * time.Second)
	}
}
