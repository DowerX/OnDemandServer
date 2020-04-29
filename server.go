package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	endtime = time.Now()
	stepsize = 20 * time.Second
	limit = 1000 * time.Second

	r := mux.NewRouter()

	r.HandleFunc("/api/reqmc", postReqMc).Methods("POST")
	http.ListenAndServe(":8000", r)
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
		//cmd := exec.Command("systemctl", "start", "papermc-server.service")
		//cmd.Start()
		lookForEnd()
	}
}

func lookForEnd() {
	looking = true
	//cmd := exec.Command("systemctl", "stop", "papermc-server.service")
	for {
		if endtime.Sub(time.Now()) <= 0 {
			//STOP
			fmt.Println("Stopping server.")
			//cmd.Start()
			looking = false
			return
		}
		time.Sleep(10 * time.Second)
	}
}
