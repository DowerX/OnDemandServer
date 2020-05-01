# OnDemandServer
Start a service by sending a POST to this server. After the time runs out the service stops. Additional request add to the countdown.

## Setup
```
git clone 
go build
```
create data.yml and config.yml

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
stepsize > time added by one request (in seconds)\
limit > max time from now to the end (in seconds)\
path > where to acces the server\
port > port of the server\
service > name of the service 

## Data
Stores user login data
