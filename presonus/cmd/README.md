# Main


## root permissions :(

We are using pcap in promiscous mode, that allows us to capture packets unrelated to this host.  Unfortunately this 
requires root permissions to run.  To start the app do:
```bash
sudo -E go run . 
```

Some explanation:
 * `sudo` runs the command as root
 * `-E` passes the current environment variables to the root user.  GOPATH and GOROOT are the most important here.
 * `go run` starts the program
 * `.` runs the app in the current folder
 
Alternatively, build the app, and then run it with sudo:
```bash
go build .
sudo ./cmd
```
  