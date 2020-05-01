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
func GetUsers() []User {
	var data, _ = ioutil.ReadFile("./data.yml")
	users := []User{}
	_ = yaml.Unmarshal(data, &users)
	return users
}

// Config _
type Config struct {
	Stepsize int
	Limit    int
	Path     string
	Port     string
	Service  string
}

// GetConfig _
func GetConfig() Config {
	var data, _ = ioutil.ReadFile("./config.yml")
	config := Config{}
	_ = yaml.Unmarshal(data, &config)
	return config
}
