package device

import (
	"context"
	"github.com/koestler/go-iotdevice/dataflow"
	"sync/atomic"
)

type Config interface {
	Name() string
	SkipFields() []string
	SkipCategories() []string
	TelemetryViaMqttClients() []string
	RealtimeViaMqttClients() []string
	LogDebug() bool
	LogComDebug() bool
}

type Device interface {
	Name() string
	Config() Config
	RegisterDb() *dataflow.RegisterDb
	IsAvailable() bool
	SubscribeAvailable(ctx context.Context) (initialAvail bool, avail <-chan bool)
	Model() string
	Run(ctx context.Context) (err error, immediateError bool)
}

type State struct {
	deviceConfig Config
	stateStorage *dataflow.ValueStorage
	registerDb   *dataflow.RegisterDb

	available        atomic.Bool
	unavailableValue dataflow.Value
	availableValue   dataflow.Value
}

func NewState(deviceConfig Config, stateStorage *dataflow.ValueStorage) State {
	registerDb := dataflow.NewRegisterDb()
	registerDb.Add(availabilityRegister)
	return State{
		deviceConfig: deviceConfig,
		stateStorage: stateStorage,
		registerDb:   registerDb,

		unavailableValue: dataflow.NewEnumRegisterValue(deviceConfig.Name(), availabilityRegister, 0),
		availableValue:   dataflow.NewEnumRegisterValue(deviceConfig.Name(), availabilityRegister, 1),
	}
}

func (c *State) Name() string {
	return c.deviceConfig.Name()
}

func (c *State) Config() Config {
	return c.deviceConfig
}

func (c *State) StateStorage() *dataflow.ValueStorage {
	return c.stateStorage
}

func (c *State) RegisterDb() *dataflow.RegisterDb {
	return c.registerDb
}

func (c *State) SetAvailable(v bool) {
	c.available.Store(v)
	if v {
		c.stateStorage.Fill(c.availableValue)
	} else {
		c.stateStorage.Fill(c.unavailableValue)
	}
}

func (c *State) IsAvailable() bool {
	return c.available.Load()
}

func (c *State) SubscribeAvailable(ctx context.Context) (initialAvail bool, availUpdate <-chan bool) {
	devName := c.Name()
	initialState, subscription := c.stateStorage.SubscribeReturnInitial(ctx, func(value dataflow.Value) bool {
		if value.DeviceName() != devName {
			return false
		}
		reg := value.Register()
		return reg.RegisterType() == dataflow.EnumRegister && reg.Name() == availabilityRegisterName
	})

	avail := c.GetAvailableByState(initialState)
	availChan := make(chan bool)
	go func() {
		defer close(availChan)
		for v := range subscription.Drain() {
			avail, updated := c.UpdateAvailable(avail, v)
			if updated {
				availChan <- avail
			}
		}
	}()

	return avail, availUpdate
}

func (c *State) GetAvailableByState(state []dataflow.Value) (avail bool) {
	devName := c.Name()
	for _, v := range state {
		if v.DeviceName() != devName {
			continue
		}
		if v.Equals(c.availableValue) {
			return true
		}
		if v.Equals(c.unavailableValue) {
			return false
		}
	}
	return false
}

func (c *State) UpdateAvailable(oldAvail bool, newValue dataflow.Value) (avail, updated bool) {
	if newValue.DeviceName() == c.Name() {
		if newValue.Equals(c.availableValue) {
			return true, true
		}
		if newValue.Equals(c.unavailableValue) {
			return false, true
		}
	}
	return oldAvail, false
}
