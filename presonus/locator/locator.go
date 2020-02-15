package locator

import (
	"net"

	"github.com/google/gopacket/pcap"
	"github.com/rltvty/go-home/logwrapper"
	"go.uber.org/zap"
)

var privateIPV4Blocks []*net.IPNet

func init() {
	privateIPV4CIDRs := []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	}

	for _, cidr := range privateIPV4CIDRs {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPV4Blocks = append(privateIPV4Blocks, block)
	}
}

func isPrivateIPV4(ip net.IP) bool {
	for _, block := range privateIPV4Blocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

//Device is a network interface
type Device struct {
	Name string
	IP   net.IP
}

//FindActiveIPV4Devices finds network devices with active, non-loopback IP4 IP addresses
func FindActiveIPV4Devices() (*[]Device, error) {
	log := logwrapper.GetInstance()
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.InfoError("Unable to get network devices", err)
		return nil, err
	}
	outDevices := []Device{}

	for _, device := range devices {
		log.Debug("Found potential device", zap.String("name", device.Name), zap.Any("info", device.Addresses))
		for _, address := range device.Addresses {
			if isPrivateIPV4(address.IP) {
				log.Debug("Found active ipv4 device", zap.String("name", device.Name), zap.String("ip", address.IP.String()))
				outDevices = append(outDevices, Device{Name: device.Name, IP: address.IP})
				break
			}
		}
	}
	return &outDevices, nil
}

func Locate() {
	FindActiveIPV4Devices()
}

/*
func OtherLocate() {

	port := 47809
	fmt.Println("Listening for UDP Broadcast Packet")

	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP: net.IPv4(0, 0, 0, 0),
		//IP:   net.IPv4( 192, 168, 1, 255 ),
		Port: port,
	})

	if err != nil {
		fmt.Println("Error listen: ", err)

	}
	for {
		data := make([]byte, 4096)
		read, remoteAddr, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("readfromudp: ", err)
		}
		for i := 0; i < read; i++ {
			fmt.Println(data[i])
		}
		fmt.Printf("Read from: %v\n", remoteAddr)
	}
}

func Locate() {
	log := logwrapper.GetInstance()
	conn, err := net.ListenPacket("udp4", ":")
	if err != nil {
		log.InfoError("Unable to open udp connection", err)
	}
	defer conn.Close()

	log.Info("Opened locator", zap.String("Address", conn.LocalAddr().String()))
	buf := make([]byte, 1024)
	for {
		deadline := time.Now().Add(30 * time.Second)
		err = conn.SetDeadline(deadline)
		n, addr, err := conn.ReadFrom(buf[0:])
		if err != nil {
			log.InfoError("error on ReadFrom:", err)
			return
		}

		//bufStr := fmt.Sprintf("0x % x\n", buf[:n])
		bufStr := fmt.Sprintf("%s", string(buf[:n]))
		log.Info("Got UDP Packet", zap.String("From Address", addr.String()), zap.String("Data", bufStr))
	}
}
*/
