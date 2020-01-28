package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jsimonetti/go-artnet/packet"
	"github.com/julienschmidt/httprouter"
	"github.com/rltvty/go-home/logwrapper"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

// The type of our middleware consists of the original handler we want to wrap and a message
type Middleware struct {
	next http.Handler
}

// Make a constructor for our middleware type since its fields are not exported (in lowercase)
func NewMiddleware(next http.Handler) *Middleware {
	return &Middleware{next: next}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Our middleware handler
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We can modify the request here; for simplicity, we will just log a message
	l := logwrapper.GetInstance()
	defer l.Sync()
	l.APIRequest(r)
	lrw := newLoggingResponseWriter(w)
	m.next.ServeHTTP(lrw, r)
	// We can modify the response here
	l.APIResponse(r, lrw.statusCode)
}

func udpAddress(ip string) string {
	return fmt.Sprintf("%s:%d", ip, packet.ArtNetPort)
}

func sendDMX(conn *net.UDPConn, node *net.UDPAddr, universe uint8, data [512]byte) {
	p := &packet.ArtDMXPacket{
		Sequence: 0,
		SubUni:   universe,
		Net:      0,
		Data:     data,
	}

	b, err := p.MarshalBinary()

	//n, err := conn.WriteTo(b, node)
	_, err = conn.WriteTo(b, node)
	if err != nil {
		fmt.Printf("error writing packet: %s\n", err)
		return
	}
	//fmt.Printf("packet sent, wrote %d bytes\n", n)
}

func getConnectedIPs() []net.IP {
	log := logwrapper.GetInstance()
	ifaces, err := net.Interfaces()
	if err != nil {
		log.PanicError("Unable to parse network interfaces", err)
	}
	ips := make([]net.IP, 0)
	for _, iface := range ifaces {
		if iface.Flags & (net.FlagLoopback | net.FlagPointToPoint) != 0 {
			continue
		}
		if !(iface.Flags & net.FlagUp != 0) {
			continue
		}

		addresses, err := iface.Addrs()
		if err != nil {
			log.InfoError("Unable to list network addresses on interface", err)
			continue
		}
		for _, address := range addresses {
			var ip net.IP
			switch v := address.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if strings.ContainsRune(ip.String(), '.') {
				log.Info("Found IP", log.String("ip", ip.String()))
				ips = append(ips, ip)
			}
		}
	}
	return ips
}

func main() {
	log := logwrapper.GetInstance()
	defer log.Sync()
	/*
		router := httprouter.New()
		router.GET("/", index)
		router.GET("/hello/:name", hello)

		log.Fatal(http.ListenAndServe(":8080", NewMiddleware(router)))
	*/

	//10.10.10.20 on universe 1 -> Sink
	//10.10.10.21 on universe 0 -> Shower

	ips := getConnectedIPs()
	if len(ips) == 0 {
		log.PanicError("No active ipv4 network interfaces found", errors.New("No interfaces found"))
	}
	ip := ips[0]

	sink, _ := net.ResolveUDPAddr("udp", udpAddress("10.10.10.20"))
	shower, _ := net.ResolveUDPAddr("udp", udpAddress("10.10.10.21"))
	src := fmt.Sprintf("%s:%d", ip.String(), packet.ArtNetPort)
	localAddr, _ := net.ResolveUDPAddr("udp", src)

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		fmt.Printf("error opening udp: %s\n", err)
		return
	}

	// set channels 1 and 7 to FL, 2-6 to zero
	// should set full red

	go func() {
		for {
			sendDMX(conn, sink, 1, [512]byte{0xFF, 0x00, 0xFF, 0x00, 0x00, 0xFF, 0xFF})
			sendDMX(conn, shower, 0, [512]byte{0x00, 0xFF, 0xFF, 0x00, 0x00, 0xFF, 0xFF})
			time.Sleep(time.Second)
		}
	}()

	for {
		time.Sleep(time.Second)
	}

}
