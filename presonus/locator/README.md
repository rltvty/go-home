# Locator

Locator makes use of a pcap library to pull UDP broadcast packets off the wire.  Docs for the library can be found here:
https://pkg.go.dev/github.com/google/gopacket

## root permissions :(

We are using pcap in promiscous mode, that allows us to capture packets unrelated to this host.  Unfortately this requires 
root permissions to run.  To start the tests do:
```bash
sudo -E go test -count=1 -v . 
```

Some explanation:
 * `sudo` runs the command as root
 * `-E` passes the current environment variables to the root user.  GOPATH and GOROOT are the most important here.
 * `go test` starts the test
 * `-count=1` prevents using cached tests results, so the test runs even if there aren't code changes
 * `-v` turns on verbose mode, so you can see the stdout logging from the process
 * `.` runs all the tests in the current folder
  
