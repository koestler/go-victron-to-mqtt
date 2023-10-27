package main

import (
	"github.com/koestler/go-iotdevice/config"
	"github.com/koestler/go-iotdevice/dataflow"
	"github.com/koestler/go-iotdevice/device"
	"github.com/koestler/go-iotdevice/httpDevice"
	"github.com/koestler/go-iotdevice/modbus"
	"github.com/koestler/go-iotdevice/modbusDevice"
	"github.com/koestler/go-iotdevice/mqttClient"
	"github.com/koestler/go-iotdevice/mqttDevice"
	"github.com/koestler/go-iotdevice/pool"
	"github.com/koestler/go-iotdevice/restarter"
	"github.com/koestler/go-iotdevice/victronDevice"
	"log"
)

func runNonMqttDevices(
	cfg *config.Config,
	modbusPool *pool.Pool[*modbus.ModbusStruct],
	stateStorage *dataflow.ValueStorage,
	commandStorage *dataflow.ValueStorage,
) (devicePool *pool.Pool[*restarter.Restarter[device.Device]]) {
	devicePool = pool.RunPool[*restarter.Restarter[device.Device]]()

	for _, deviceConfig := range cfg.VictronDevices() {
		if cfg.LogWorkerStart() {
			log.Printf("device[%s]: start victron type", deviceConfig.Name())
		}

		deviceConfig := victronDeviceConfig{deviceConfig}
		dev := victronDevice.NewDevice(deviceConfig, deviceConfig, stateStorage)
		watchedDev := restarter.CreateRestarter[device.Device](deviceConfig, dev)
		watchedDev.Run()
		devicePool.Add(watchedDev)
	}

	for _, deviceConfig := range cfg.ModbusDevices() {
		if cfg.LogWorkerStart() {
			log.Printf("device[%s]: start modbus type", deviceConfig.Name())
		}

		deviceConfig := modbusDeviceConfig{deviceConfig}
		modbusInstance := modbusPool.GetByName(deviceConfig.Bus())
		if modbusInstance == nil {
			log.Printf("device[%s]: start failed: bus=%s unavailable", deviceConfig.Name(), deviceConfig.Bus())
			continue
		}

		dev := modbusDevice.NewDevice(deviceConfig, deviceConfig, modbusInstance, stateStorage, commandStorage)
		watchedDev := restarter.CreateRestarter[device.Device](deviceConfig, dev)
		watchedDev.Run()
		devicePool.Add(watchedDev)
	}

	for _, deviceConfig := range cfg.HttpDevices() {
		if cfg.LogWorkerStart() {
			log.Printf("device[%s]: start tearacom type", deviceConfig.Name())
		}

		deviceConfig := httpDeviceConfig{deviceConfig}
		dev := httpDevice.NewDevice(deviceConfig, deviceConfig, stateStorage, commandStorage)
		watchedDev := restarter.CreateRestarter[device.Device](deviceConfig, dev)
		watchedDev.Run()
		devicePool.Add(watchedDev)
	}

	return
}

func runMqttDevices(
	cfg *config.Config,
	devicePool *pool.Pool[*restarter.Restarter[device.Device]],
	mqttClientPool *pool.Pool[mqttClient.Client],
	stateStorage *dataflow.ValueStorage,
) {
	for _, deviceConfig := range cfg.MqttDevices() {
		if cfg.LogWorkerStart() {
			log.Printf("device[%s]: start mqtt type", deviceConfig.Name())
		}

		deviceConfig := mqttDeviceConfig{deviceConfig}
		dev := mqttDevice.NewDevice(deviceConfig, deviceConfig, stateStorage, mqttClientPool)
		watchedDev := restarter.CreateRestarter[device.Device](deviceConfig, dev)
		watchedDev.Run()
		devicePool.Add(watchedDev)
	}
}

// the following structs / methods are used to cast config.RegisterFilterConfig into dataflow.RegisterFilterConf

type victronDeviceConfig struct {
	config.VictronDeviceConfig
}

func (c victronDeviceConfig) RegisterFilter() dataflow.RegisterFilterConf {
	return c.RegisterFilter()
}

type modbusDeviceConfig struct {
	config.ModbusDeviceConfig
}

func (c modbusDeviceConfig) RegisterFilter() dataflow.RegisterFilterConf {
	return c.RegisterFilter()
}

type httpDeviceConfig struct {
	config.HttpDeviceConfig
}

func (c httpDeviceConfig) RegisterFilter() dataflow.RegisterFilterConf {
	return c.RegisterFilter()
}

type mqttDeviceConfig struct {
	config.MqttDeviceConfig
}

func (c mqttDeviceConfig) RegisterFilter() dataflow.RegisterFilterConf {
	return c.RegisterFilter()
}
