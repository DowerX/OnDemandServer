package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"./data"
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
	c = data.GetConfig()
	endtime = time.Now()
	stepsize = time.Duration(c.Stepsize) * time.Second
	limit = time.Duration(c.Limit) * time.Second

	lf, err := os.OpenFile("./log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("ERR: Can't log!")
		return
	}
	log.SetOutput(lf)

	r := mux.NewRouter()
	r.HandleFunc(c.Path, postReqMc).Methods("POST")
	http.ListenAndServe(c.Port, r)
}

func postReqMc(w http.ResponseWriter, r *http.Request) {

	creds := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println("WARN: Incoming request, but wrong format!")
		fmt.Println("Incoming request, but wrong format!")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users := data.GetUsers()

	for i := 0; i < len(users); i++ {
		if users[i].Username == creds.Username {
			c := bcrypt.CompareHashAndPassword([]byte(users[i].PasswordHash), []byte(creds.Password))
			if c == nil {
				go req(users[i])
				log.Println("INFO: Incoming request: ", creds.Username)
				fmt.Println("INFO: Incoming request: ", creds.Username)
				w.WriteHeader(http.StatusAccepted)
				return
			}
		}
	}
	log.Println("WARN: Unauthorized request: ", creds.Username)
	fmt.Println("WARN: Unauthorized request: ", creds.Username)
	w.WriteHeader(http.StatusUnauthorized)
}

func req(user data.User) {
	if endtime.Sub(time.Now()) <= 0 {
		endtime = time.Now()
	}

	if endtime.Sub(time.Now()) <= limit {
		endtime = endtime.Add(stepsize)
	}

	fmt.Println("INFO: Added", stepsize, "seconds.", "Service will stop at ", endtime, ".")
	log.Println("INFO: Added", stepsize, "seconds.", "Service will stop at ", endtime, ".")

	if looking == false {
		//START
		fmt.Println("Starting service.")
		cmd := exec.Command("systemctl", "start", c.Service)
		cmd.Start()
		log.Println("INFO: Started service.")
		lookForEnd()
	}
}

func lookForEnd() {
	looking = true
	cmd := exec.Command("systemctl", "stop", c.Service)
	for {
		if endtime.Sub(time.Now()) <= 0 {
			//STOP
			fmt.Println("Stopping server.")
			cmd.Start()
			log.Println("INFO: Stopped service.")
			looking = false
			return
		}
		time.Sleep(10 * time.Second)
	}
}
