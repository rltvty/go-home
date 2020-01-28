package netutils

import (
	"net"
	"strings"

	"github.com/rltvty/go-home/logwrapper"
)

//GetConnectedIPV4s returns a list of IPV4 IPs currently active on the host
func GetConnectedIPV4s() []net.IP {
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