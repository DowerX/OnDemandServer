package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func main() {
	c := data.GetConfig()
	endtime = time.Now()
	stepsize = time.Duration(c.Stepsize) * time.Second
	limit = time.Duration(c.Limit) * time.Second

	r := mux.NewRouter()
	r.HandleFunc(c.Path, postReqMc).Methods("POST")
	http.ListenAndServe(c.Port, r)
}

func postReqMc(w http.ResponseWriter, r *http.Request) {
	creds := Credentials{}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users := data.GetUsers()

	for i := 0; i < len(users); i++ {
		if users[i].Username == creds.Username {
			c := bcrypt.CompareHashAndPassword([]byte(users[i].PasswordHash), []byte(creds.Password))
			if c == nil {
				go req(users[i])
				w.WriteHeader(http.StatusAccepted)
				return
			}
		}
	}
	w.WriteHeader(http.StatusUnauthorized)
}

func req(user data.User) {
	fmt.Println(user.Username, "requested the server.")
	if endtime.Sub(time.Now()) <= 0 {
		endtime = time.Now()
	}

	if endtime.Sub(time.Now()) <= limit {
		endtime = endtime.Add(stepsize)
	}

	fmt.Println("Added ", stepsize, ", now at ", endtime, ".")

	if looking == false {
		//START
		fmt.Println("Starting server.")
		cmd := exec.Command("systemctl", "start", c.Service)
		cmd.Start()
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
			looking = false
			return
		}
		time.Sleep(10 * time.Second)
	}
}
