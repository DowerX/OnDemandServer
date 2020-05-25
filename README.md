# OnDemandServer
Start a service by sending a POST to this server. After the time runs out the service stops. Additional request add to the countdown.

## Setup
```
git clone / go get
go build
```
set up user credentials (data.yml by default) 

## Usage
Send POST from your website or with Postman.\
The body should look like this:
```
{
	"username" : "USERNAME",
	"password" : "PASSWORD"
}
```

## Config
config > set config file\
stepsize > time added by one request (in seconds)\
limit > max time from now to the end (in seconds)\
path > where to acces the server\
port > port of the server\
startscript > path to the startig script\
stopscript > path to the stopping script\
log > enable/disable logging to a file\
logfile > path to the logfile\
users > path to the users file

## Data
Stores user login data
