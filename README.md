# Network-CLI

## Simple Local Network Command Line Interface

### Get Started:

First to build and run the program, change to the project directory and run the following command. This output the binary executable called **network**

``` Bash
$ go build -o network
```
Note: Make sure your machine has [Go](https://golang.org/dl/) environment.

Once you generated the executable file, put **network** to your _/usr/bin/_ directory.

### Usages:

- Server:

``` Bash
$ network server -port 3000         # create localhost on port 3000 (default 8080)
$ network server -port 3000 -file   # create localhost on port 3000 and use files as server content
```

- Forward:

``` Bash
$ network forward -target 127.0.0.1:3000 -port 8080   # forward 127.0.0.1 port 3000 to local port 8080
```

- Check

``` Bash
$ network check                               # Return all of the unavailable port number on local machine
$ network check -portList 80,1024,3000,8000   # Return all of the unavailable port number on local machine from given port list
$ network check -ip                           # Return local machine's internal and external IP addresses
```

---

More command details use `$ network help`
