package main

import (
	"github.com/rltvty/go-home/logwrapper"
	"github.com/rltvty/go-home/presonus/locator"
	"github.com/rltvty/go-home/presonus/simple_conn"
	"go.uber.org/zap"
)

func main()  {
	log := logwrapper.GetInstance()
	events := make(chan locator.PresonusDeviceEvent)

	go func() {
		devices := make(map [*locator.PresonusDevice]*simple_conn.Client)
		for event := range events {
			log.Info("Presonus Device Event", zap.Bool("isAdd", event.IsAdd), zap.Any("device", event.Device))
			if event.IsAdd {
				device := simple_conn.Device{
					Kind: event.Device.Kind,
					IP:   event.Device.IP.String(),
					Port: event.Device.Port,
				}
				client := simple_conn.InitClient(device)
				devices[&event.Device] = client
				go func() {
					select {
						case read, ok := <- client.Read:
							if !ok {
								log.Info("returning after client 'read' channel has closed")
								return
							}
							log.Debug("reading from client channel", zap.ByteString("message", read))
					}
				}()
			} else {
				client, found := devices[&event.Device]
				if found {
					client.Close()
					delete(devices, &event.Device)
				}
			}
		}
	}()

	locator.MainLoop(events)
}
