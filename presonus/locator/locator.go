package locator

import (
	"encoding/binary"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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

	err = handle.SetBPFFilter("ip broadcast")
	if err != nil {
		log.InfoError("Error setting pcap filter", err)
	}

	// Start processing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetCount := 0
	for packet := range packetSource.Packets() {
		// Process packet here
		log.Debug("")
		log.Debug("")
		log.Debug("Got Packet", zap.Any("Packet", packet))

		packetCount++

		// Iterate over all layers, printing out each layer type
		for _, layer := range packet.Layers() {
			log.Debug("PACKET LAYER:", zap.Any("LayerType", layer.LayerType()))
		}

		// Get the Ethernet layer from this packet
		if ethernetLayer := packet.Layer(layers.LayerTypeEthernet); ethernetLayer != nil {
			// Get actual Ethernet data from this layer
			ethernet, _ := ethernetLayer.(*layers.Ethernet)
			log.Debug("MAC", zap.Any("src", ethernet.SrcMAC), zap.Any("dst", ethernet.DstMAC))
		}

		// Get the IPv4 layer from this packet
		if ipv4Layer := packet.Layer(layers.LayerTypeIPv4); ipv4Layer != nil {
			// Get actual IPv4 data from this layer
			ipv4, _ := ipv4Layer.(*layers.IPv4)
			log.Debug("IP", zap.Any("src", ipv4.SrcIP), zap.Any("dst", ipv4.DstIP))
		}

		// Get the UDP layer from this packet
		if udpLayer := packet.Layer(layers.LayerTypeUDP); udpLayer != nil {
			// Get actual UDP data from this layer
			udp, _ := udpLayer.(*layers.UDP)
			log.Debug("Ports", zap.Any("src", udp.SrcPort), zap.Any("dst", udp.DstPort))
		}

		// Get the Application layer from this packet
		if app := packet.ApplicationLayer(); app != nil {
			log.Debug("Payload", zap.Any("data", string(app.Payload())))
			item, _ := DecodeData(app.Payload(), "locator")
			if item != nil {
				log.Debug("Found item", zap.Any(item.Kind, item))
			}
		}

		// Only capture 100 and then stop
		if packetCount > 100 {
			break
		}
	}
}

type Data struct {
	Source string
	Mode string
	Port uint16
	Model string
	MacAddress string
	Kind string
}

func DecodeData(payload []byte, source string) (*Data, error) {
	log := logwrapper.GetInstance()
	if string(payload[0:2]) == "UC" {
		code := binary.LittleEndian.Uint16(payload[6:])
		switch code {
		case 16708:
			log.Debug("Found Presounus Broadcast")
			if len(payload) <= 50 {
				return &Data{
					Source: source,
					Mode: "broadcast",
					Port: binary.LittleEndian.Uint16(payload[4:]),
					Model: strings.TrimSpace(strings.ReplaceAll(string(payload[16:28]), string(0), " ")),
					MacAddress: string(payload[28:45]),
					Kind: "speaker",
				}, nil
			} else {
				return &Data{
					Source: source,
					Mode: "broadcast",
					Port: binary.LittleEndian.Uint16(payload[4:]),
					Model: strings.TrimSpace(strings.ReplaceAll(string(payload[32:50]), string(0), " ")),
					MacAddress: "",
					Kind: "mixer",
				}, nil
			}
		}
	}
	return nil, nil
}
