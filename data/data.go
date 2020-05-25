package data

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// User _
type User struct {
	PasswordHash string
	Username     string
}

// GetUsers _
func GetUsers(f string) []User {
	var data, err = ioutil.ReadFile(f)
	users := []User{}
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &users)
	if err != nil {
		panic(err)
	}
	return users
}

// Config _
type Config struct {
	Stepsize    int
	Limit       int
	Path        string
	Port        string
	StartScript string
	StopScript  string
	Log         bool
	Logfile     string
	Users       string
}

// GetConfig _
func GetConfig(f string) Config {
	var data, err = ioutil.ReadFile(f)
	config := Config{}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}
