package simple_conn
// based on connection.go, but simplified to just be a connection and not also a connection manager.

import (
	"encoding/hex"
	"fmt"
	"github.com/rltvty/go-home/logwrapper"
	"go.uber.org/zap"
	"net"
	"time"
)

const initRequest = "554300010800554d00006400cfde"

//'UC§JMfd{"id": "Subscribe","clientName": "SL Room Control","clientType": "Mac","clientDescription": "Eric\'s MacBook Air","clientIdentifier": "Ericâs MacBook Air"}UCKAfd'
const subscriptionRequest = "55430001a7004a4d660064009d0000007b226964223a2022537562736372696265222c22636c69656e744e616d65223a2022534c20526f6f6d20436f6e74726f6c222c22636c69656e7454797065223a20224d6163222c22636c69656e744465736372697074696f6e223a2022457269635c2773204d6163426f6f6b20416972222c22636c69656e744964656e746966696572223a202245726963e2809973204d6163426f6f6b20416972227d5543000106004b4166006400"

//UCKAed
const keepAlive = "5543000106004b4166006400"

type Device struct {
	Kind string
	IP   string
	Port uint16
}

type Client struct {
	socket net.Conn
	write  chan []byte
	Read   chan []byte
	isOpen bool
}

func InitClient(device Device) *Client {
	log := logwrapper.GetInstance()
	address := fmt.Sprintf("%s:%v", device.IP, device.Port)
	connection, err := net.Dial("tcp", address)
	if err != nil {
		log.InfoError("error opening socket connection", err)
		return nil
	}
	client := &Client{socket: connection, write: make(chan []byte), Read: make(chan []byte), isOpen: true}

	go client.receive()
	go client.send()
	go client.keepAlive()

	err = client.WriteHex(initRequest)
	if err != nil {
		log.InfoError("error sending init request", err)
		return nil
	}
	err = client.WriteHex(subscriptionRequest)
	if err != nil {
		log.InfoError("error sending subscription request", err)
		return nil
	}

	return client
}

func (client *Client) receive() {
	log := logwrapper.GetInstance()
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			log.InfoError("could not read from socket", err)
			client.Close()
			return
		}
		if length > 0 {
			client.Read <- message
			log.Debug("client received data", zap.Int("length", length), zap.ByteString("message", message))
		}
	}
}

func (client *Client) send() {
	log := logwrapper.GetInstance()
	defer client.Close()
	for {
		select {
		case message, ok := <-client.write:
			if !ok {
				log.Info("returning after 'write' channel has closed")
				return
			}
			_, err := client.socket.Write(message)
			if err != nil {
				log.InfoError("could not write to socket", err)
				client.Close()
				return
			}
		}
	}
}

func (client *Client) keepAlive() {
	for {
		time.Sleep(time.Second * 3)
		err := client.WriteHex(keepAlive)
		if err != nil {
			log := logwrapper.GetInstance()
			log.InfoError("keep alive ending", err)
			return
		}
	}
}

func (client *Client) Close() {
	log := logwrapper.GetInstance()
	if client.isOpen {
		client.isOpen = false
		err := client.socket.Close()
		if err != nil {
			log.InfoError("error closing socket", err)
		}
		close(client.write)
		close(client.Read)
	}
}

func (client *Client) WriteHex(hexString string) error {
	log := logwrapper.GetInstance()
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		log.InfoError("error decoding hex string", err)
	}
	if client.isOpen {
		client.write <- bytes
		return nil
	} else {
		return fmt.Errorf("connection is closed")
	}
}
