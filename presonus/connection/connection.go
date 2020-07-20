package connection

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"
)

const init_request = "554300010800554d00006400cfde"

//'UC§JMfd{"id": "Subscribe","clientName": "SL Room Control","clientType": "Mac","clientDescription": "Eric\'s MacBook Air","clientIdentifier": "Ericâs MacBook Air"}UCKAfd'
const subscription_request = "55430001a7004a4d660064009d0000007b226964223a2022537562736372696265222c22636c69656e744e616d65223a2022534c20526f6f6d20436f6e74726f6c222c22636c69656e7454797065223a20224d6163222c22636c69656e744465736372697074696f6e223a2022457269635c2773204d6163426f6f6b20416972222c22636c69656e744964656e746966696572223a202245726963e2809973204d6163426f6f6b20416972227d5543000106004b4166006400"

//UCKAed
const keep_alive = "5543000106004b4166006400"

type Device struct {
	Kind string
	IP   string
	Port uint16
}

/*
func Connect(device Device) {
	log := logwrapper.GetInstance()
	client, err := net.Dial("tcp", fmt.Sprintf("%s:%v", device.IP, device.Port))
	if err != nil {
		log.InfoError(fmt.Sprintf("unable to dial %s", device.Kind), err)
	}

}*/

type ClientManager struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket net.Conn
	data   chan []byte
}

func (manager *ClientManager) start() {
	for {
		select {
		case connection := <-manager.register:
			manager.clients[connection] = true
			fmt.Println("Added new connection!")
		case connection := <-manager.unregister:
			if _, ok := manager.clients[connection]; ok {
				close(connection.data)
				delete(manager.clients, connection)
				fmt.Println("A connection has terminated!")
			}
		}
	}
}

func (manager *ClientManager) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)
		if err != nil {
			fmt.Println("Manager Error Message:")
			fmt.Println(err)
			manager.unregister <- client
			client.socket.Close()
			break
		}
		if length > 0 {
			fmt.Printf("Manager RECEIVED %v bytes, contents %s\n", length, message)
		}
	}
}

func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				fmt.Println("returning after not getting data from client data socket")
				return
			}
			client.socket.Write(message)
		}
	}
}

var manager *ClientManager

func StartManager() {
	fmt.Println("Starting manager...")
	manager = &ClientManager{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go manager.start()
}

func InitConnection(device Device) {
	address := fmt.Sprintf("%s,%v", device.IP, device.Port)
	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
	}
	client := &Client{socket: connection, data: make(chan []byte)}

	manager.register <- client
	go manager.receive(client)
	go manager.send(client)

	writeHex(client, init_request)
	writeHex(client, subscription_request)

	for {
		time.Sleep(time.Second * 3)
		writeHex(client, keep_alive)
	}
}

func writeHex(client *Client, hexString string) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println(err)
	}
	client.data <- bytes
}
