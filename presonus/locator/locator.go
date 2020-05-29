package locator

import (
	"github.com/google/gopacket"
	"net"
	"strings"

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
	log := logwrapper.GetInstance()
	log.Info("starting")
	devices, err :=	FindActiveIPV4Devices()
	if err != nil {
		log.InfoError("Unable to find Active IP4 devices", err)
	}
	log.Info("Found IP4 Devices: ", zap.Any("devices", devices))

	var chosenDevice Device
	for _, device := range *devices {
		if strings.HasPrefix(device.Name, "en") {
			chosenDevice = device
		}
	}
	log.Info("Using Device: ", zap.String("name", chosenDevice.Name))

	inactive, err := pcap.NewInactiveHandle(chosenDevice.Name)
	if err != nil {
		log.InfoError("Unable to create pcap handle", err)
	}
	defer inactive.CleanUp()

	// Call various functions on inactive to set it up the way you'd like:
	if err = inactive.SetTimeout(pcap.BlockForever); err != nil {
		log.InfoError("Error setting pcap timeout", err)
	} else if err = inactive.SetPromisc(true); err != nil {
		log.InfoError("Error setting pcap promiscuous mode", err)
	}

	// Finally, create the actual handle by calling Activate:
	handle, err := inactive.Activate()  // after this, inactive is no longer valid
	if err != nil {
		log.InfoError("Error activating pcap handle", err)
	}
	defer handle.Close()

	// Start processing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetCount := 0
	for packet := range packetSource.Packets() {
		// Process packet here
		log.Debug("Got Packet", zap.Any("Packet", packet))

		packetCount++

		// Only capture 100 and then stop
		if packetCount > 100 {
			break
		}
	}
}
