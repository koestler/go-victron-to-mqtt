package mqttClient

import (
	"github.com/koestler/go-iotdevice/dataflow"
	"github.com/koestler/go-iotdevice/device"
	"log"
	"sync"
)

type ClientStruct struct {
	cfg                Config
	devicePoolInstance *device.DevicePool
	storage            *dataflow.ValueStorageInstance

	shutdown chan struct{}

	subscriptionsMutex sync.RWMutex
	subscriptions      []subscription
}

type subscription struct {
	subscribeTopic string
	messageHandler MessageHandler
}

func createClientStruct(
	cfg Config,
	devicePoolInstance *device.DevicePool,
	storage *dataflow.ValueStorageInstance,
) ClientStruct {
	return ClientStruct{
		cfg:                cfg,
		devicePoolInstance: devicePoolInstance,
		storage:            storage,
		shutdown:           make(chan struct{}),
	}
}

func (c *ClientStruct) AddRoute(subscribeTopic string, messageHandler MessageHandler) {
	log.Printf("mqttClient[%s]: add route for topic='%s'", c.cfg.Name(), subscribeTopic)

	s := subscription{subscribeTopic: subscribeTopic}

	if c.cfg.LogMessages() {
		s.messageHandler = func(message Message) {
			// only log first 80 chars of payload
			pl := make([]byte, 0, 80)
			pl = append(pl, message.Payload()[:80]...)
			if len(message.Payload()) > 80 {
				pl = append(pl, []byte("...")...)
			}

			log.Printf("mqttClient[%s]: received: %s %s", c.cfg.Name(), message.Topic(), pl)
			messageHandler(message)
		}
	} else {
		s.messageHandler = func(message Message) {
			messageHandler(message)
		}
	}

	c.subscriptionsMutex.Lock()
	defer c.subscriptionsMutex.Unlock()
	c.subscriptions = append(c.subscriptions, s)
}

func (c *ClientStruct) Name() string {
	return c.cfg.Name()
}
