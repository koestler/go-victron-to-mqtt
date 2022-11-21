package mqttClient

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/koestler/go-iotdevice/dataflow"
	"github.com/koestler/go-iotdevice/device"
	"log"
	"net/url"
	"time"
)

type ClientV5 struct {
	ClientStruct

	cliCfg autopaho.ClientConfig
	cm     *autopaho.ConnectionManager
	router *paho.StandardRouter

	ctx    context.Context
	cancel context.CancelFunc
}

func CreateV5(
	cfg Config,
	devicePoolInstance *device.DevicePool,
	storage *dataflow.ValueStorageInstance,
) (client *ClientV5) {
	ctx, cancel := context.WithCancel(context.Background())
	client = &ClientV5{
		ClientStruct: createClientStruct(cfg, devicePoolInstance, storage),
		router:       paho.NewStandardRouter(),
		ctx:          ctx,
		cancel:       cancel,
	}

	// configure mqtt library
	client.cliCfg = autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{cfg.Broker()},
		KeepAlive:         uint16(cfg.KeepAlive().Seconds()),
		ConnectRetryDelay: cfg.ConnectRetryDelay(),
		ConnectTimeout:    cfg.ConnectTimeout(),
		OnConnectionUp:    client.onConnectionUp(),
		OnConnectError: func(err error) {
			log.Printf("mqttClientV5[%s]: connection error: %s", cfg.Name(), err)
		},
		ClientConfig: paho.ClientConfig{
			ClientID: cfg.ClientId(),
			Router:   client.router,
		},
	}

	// setup logging
	if cfg.LogDebug() {
		prefix := fmt.Sprintf("mqttClientV5[%s]: ", cfg.Name())
		client.cliCfg.Debug = logger{prefix: prefix + "autoPaho: "}
		client.cliCfg.PahoDebug = logger{prefix: prefix + "paho: "}
	}

	// configure login
	if user := cfg.User(); len(user) > 0 {
		client.cliCfg.SetUsernamePassword(user, []byte(cfg.Password()))
	}

	// setup availability topic using will
	if client.AvailabilityEnabled() {
		client.cliCfg.SetWillMessage(client.GetAvailabilityTopic(), []byte(availabilityOffline), cfg.Qos(), true)
	}

	return
}

func (c *ClientV5) Run() {
	// add routes to router
	c.subscriptionsMutex.RLock()
	defer c.subscriptionsMutex.RUnlock()
	for _, s := range c.subscriptions {
		sub := s
		c.router.RegisterHandler(sub.subscribeTopic, func(p *paho.Publish) {
			sub.messageHandler(Message{
				topic:   p.Topic,
				payload: p.Payload,
			})
		})
	}

	// start connection manager
	var err error
	c.cm, err = autopaho.NewConnection(c.ctx, c.cliCfg)
	if err != nil {
		panic(err) // never happens
	}
}

func (c *ClientV5) onConnectionUp() func(*autopaho.ConnectionManager, *paho.Connack) {
	return func(cm *autopaho.ConnectionManager, conack *paho.Connack) {
		log.Printf("mqttClientV5[%s]: connection is up", c.cfg.Name())
		// publish online
		if c.AvailabilityEnabled() {
			go func() {
				_, err := cm.Publish(c.ctx, c.availabilityMsg(availabilityOnline))
				if err != nil {
					log.Printf("mqttClientV5[%s]: error during publish: %s", c.cfg.Name(), err)
				}
			}()
		}
		// subscribe topics
		if _, err := cm.Subscribe(c.ctx, &paho.Subscribe{
			Subscriptions: func() (ret map[string]paho.SubscribeOptions) {
				c.subscriptionsMutex.RLock()
				defer c.subscriptionsMutex.RUnlock()
				ret = make(map[string]paho.SubscribeOptions, len(c.subscriptions))

				subOpts := paho.SubscribeOptions{QoS: c.cfg.Qos()}
				for _, s := range c.subscriptions {
					ret[s.subscribeTopic] = subOpts
				}
				return
			}(),
		}); err != nil {
			log.Printf("mqttClientV5[%s]: failed to subscribe: %s", c.cfg.Name(), err)
		}

		// setup Realtime (send data as soon as it arrives) output
		if c.cfg.RealtimeEnable() {
			// transmitRealtime values from data store and publish to mqtt broker
			go func() {
				// setup empty filter (everything)
				subscription := c.storage.Subscribe(dataflow.Filter{})
				defer subscription.Shutdown()
				for {
					select {
					case <-c.ctx.Done():
						return
					case value := <-subscription.GetOutput():
						if p, err := c.getRealtimePublishMessage(value); err != nil {
							if _, err := c.cm.Publish(c.ctx, p); err != nil {
								log.Printf("mqttClientV5[%s]: cannot publish realtime: %s", c.cfg.Name(), err)
							}
						}
					}
				}
			}()
			log.Printf("mqttClient[%s]: start sending realtime stat messages", c.cfg.Name())
		}

		// setup Telemetry support
		if interval := c.cfg.TelemetryInterval(); interval > 0 {
			go func() {
				ticker := time.NewTicker(interval)
				for {
					select {
					case <-c.ctx.Done():
						return
					case <-ticker.C:
						for deviceName, dev := range c.devicePoolInstance.GetDevices() {
							deviceFilter := dataflow.Filter{IncludeDevices: map[string]bool{deviceName: true}}
							values := c.storage.GetSlice(deviceFilter)
							if p, err := c.getTelemetryPublishMessage(deviceName, dev, values); err != nil {
								if _, err := c.cm.Publish(c.ctx, p); err != nil {
									log.Printf("mqttClientV5[%s]: cannot publish telemetry: %s", c.cfg.Name(), err)
								}
							}

						}
					}
				}
			}()

			log.Printf("mqttClientV5[%s]: start sending telemetry messages every %s", c.cfg.Name(), interval.String())
		}

	}
}

func (c *ClientV5) Shutdown() {
	close(c.shutdown)

	// publish availability offline
	if c.AvailabilityEnabled() {
		ctx, cancel := context.WithTimeout(c.ctx, time.Second)
		defer cancel()
		_, err := c.cm.Publish(ctx, c.availabilityMsg(availabilityOffline))
		if err != nil {
			log.Printf("mqttClientV5[%s]: error during publish: %s", c.cfg.Name(), err)
		}
	}

	ctx, cancel := context.WithTimeout(c.ctx, time.Second)
	defer cancel()
	if err := c.cm.Disconnect(ctx); err != nil {
		log.Printf("mqttClientV5[%s]: error during disconnect: %s", c.cfg.Name(), err)
	}

	// cancel main context
	c.cancel()

	log.Printf("mqttClientV5[%s]: shutdown completed", c.cfg.Name())
}

func (c *ClientV5) availabilityMsg(payload string) *paho.Publish {
	return &paho.Publish{
		QoS:     c.cfg.Qos(),
		Topic:   c.GetAvailabilityTopic(),
		Payload: []byte(payload),
		Retain:  availabilityRetain,
	}
}
