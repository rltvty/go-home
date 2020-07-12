package locator_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/rltvty/go-home/logwrapper"
	"go.uber.org/zap"

	//. "github.com/onsi/gomega"

	. "github.com/rltvty/go-home/presonus/locator"
)

var _ = Describe("Locator", func() {
	Describe("MainLoop", func() {
		log := logwrapper.GetInstance()
		events := make(chan PresonusDeviceEvent)
		go func() {
			for event := range events {
				log.Info("Presonus Device Event", zap.Any(event.EventType, event.Device))
			}
		}()

		MainLoop(events)
	})
})
