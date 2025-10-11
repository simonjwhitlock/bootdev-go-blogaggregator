# bootdev-go-blogaggregator

This is a simple CLI based RSS feed agrigator. there is no fancy UI or other bits to this.

## Requirements:

This requires access to a postgres instance (may be local or remote).

The host system requires GO installed.

## Configuration:
### Define the programs congfiguration
The Configureation is simple. in the root of the system users 'home' folder (C:/Users/%username% on widnows), Create a file called *.gatorconfig.json* in this file enter the following:

```
{"db_url":"postgres://<dbUsername>:<dbPW>@<DBhostname/IP>:<dbhostport>/<dbname>"}
```
you may need to add ```?sslmode=disable``` to the end of the db connection string if your hosted locally or on a sever without SSL setup.

### Initalising the data base:
install the go module 'goose' using
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
in the terminal navagate to *sql/schema/* inside the project repo then implimenmt the database schema with:
```
goose postgres "<db connection sring as above>" up
```
this will upgrade the database to v5 (the latest at this time)

return back to the root of the repo.

## Commands
all commands are carried out from a terminal navagated to the root of the repo. the availabe commands are:

- login - log into an existing user,  ```go run . login <username>``` 
- register - create a new user, atomatically changes to this new user ```go run . register <username>```
- reset - clear all data out of database (incuding users) ```go run . reset```
- users - list all registered users ```go run . users```
- agg - constantly gets posts in all feeds at intervals, ```go run . agg <interval eg: 1m>```
- addfeed - add a new feed (auto follows for logged in user), ```go run . addfeed <feedname> <url>"```
- feeds - list feeds already in system, ```go run . feeds```
- follow - follow an existing feed, ```go run . follow <url>```
- following - list feeds user is following ```go run . following```
- unfollow - unfollow a feed,  ```go run . unfollow <url>```
- browse - show the latest posts in users feed ```go run . browse <number of posts (optional, defaults to 2)>```

## Useage

Once configured (database setup and .gatorconfig.json has database connection in it), the first step is to register a user and add an RSS feed to follow.

create user:
```
go run . register <username>
```
register a feed:
```
go run . addfeed <feedname> <feeddURL>
```
In a new terminal start the Agrigator function, this takes a time interval as per the otions in https://pkg.go.dev/time#ParseDuration : 
```
go run . agg <timeinterval>
```
Return to original terminal for all other interations

list the latest published post in users feed:
```
go run . browse <number of posts to get>
```

### Add another usr to the system
this allows multiple users, these are not secured at all, just used to curate separate lists.

to add new user and change to it
```
go run . register <new username>
```

### If a feed exists in the database but is not followed by the current user:
list all existing feeds
```
go run . feeds
```
to follow an existing feed:
```
go run . follow <feed url>
```

### Change logged on user
list all users
```
go run . users
```
change user
```
go run . logon <username>
```
