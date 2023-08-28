package config

import (
	"fmt"
	"golang.org/x/exp/maps"
)

func (c Config) MarshalYAML() (interface{}, error) {
	return configRead{
		Version:         &c.version,
		ProjectTitle:    c.projectTitle,
		LogConfig:       &c.logConfig,
		LogWorkerStart:  &c.logWorkerStart,
		LogStorageDebug: &c.logStorageDebug,
		HttpServer:      convertEnableableToRead[HttpServerConfig, httpServerConfigRead](c.httpServer),
		Authentication:  convertEnableableToRead[AuthenticationConfig, authenticationConfigRead](c.authentication),
		MqttClients:     convertMapToRead[MqttClientConfig, mqttClientConfigRead](c.mqttClients),
		Modbus:          convertMapToRead[ModbusConfig, modbusConfigRead](c.modbus),
		VictronDevices:  convertMapToRead[VictronDeviceConfig, victronDeviceConfigRead](c.victronDevices),
		ModbusDevices:   convertMapToRead[ModbusDeviceConfig, modbusDeviceConfigRead](c.modbusDevices),
		HttpDevices:     convertMapToRead[HttpDeviceConfig, httpDeviceConfigRead](c.httpDevices),
		MqttDevices:     convertMapToRead[MqttDeviceConfig, mqttDeviceConfigRead](c.mqttDevices),
		Views:           convertListToRead[ViewConfig, viewConfigRead](c.views),
		HassDiscovery:   convertListToRead[HassDiscovery, hassDiscoveryRead](c.hassDiscovery),
	}, nil
}

type convertable[O any] interface {
	convertToRead() O
}

type enableable[O any] interface {
	Enabled() bool
	convertable[O]
}

func convertEnableableToRead[I enableable[O], O any](inp I) *O {
	if !inp.Enabled() {
		return nil
	}
	r := inp.convertToRead()
	return &r
}

type mappable[O any] interface {
	Nameable
	convertable[O]
}

func convertMapToRead[I mappable[O], O any](inp []I) (oup map[string]O) {
	oup = make(map[string]O, len(inp))
	for _, c := range inp {
		oup[c.Name()] = c.convertToRead()
	}
	return
}

func convertListToRead[I convertable[O], O any](inp []I) (oup []O) {
	oup = make([]O, len(inp))
	i := 0
	for _, c := range inp {
		oup[i] = c.convertToRead()
		i++
	}
	return
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c HttpServerConfig) convertToRead() httpServerConfigRead {
	frontendProxy := ""
	if c.frontendProxy != nil {
		frontendProxy = c.frontendProxy.String()
	}

	return httpServerConfigRead{
		Bind:            c.bind,
		Port:            &c.port,
		LogRequests:     &c.logRequests,
		FrontendProxy:   frontendProxy,
		FrontendPath:    c.frontendPath,
		FrontendExpires: c.frontendExpires.String(),
		ConfigExpires:   c.configExpires.String(),
		LogDebug:        &c.logDebug,
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c AuthenticationConfig) convertToRead() authenticationConfigRead {
	jwtSecret := string(c.jwtSecret)
	return authenticationConfigRead{
		JwtSecret:         &jwtSecret,
		JwtValidityPeriod: c.jwtValidityPeriod.String(),
		HtaccessFile:      &c.htaccessFile,
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c MqttClientConfig) convertToRead() mqttClientConfigRead {
	return mqttClientConfigRead{
		Broker:            c.broker.String(),
		ProtocolVersion:   &c.protocolVersion,
		User:              c.user,
		Password:          c.password,
		ClientId:          &c.clientId,
		Qos:               &c.qos,
		KeepAlive:         c.keepAlive.String(),
		ConnectRetryDelay: c.connectRetryDelay.String(),
		ConnectTimeout:    c.connectTimeout.String(),
		AvailabilityTopic: &c.availabilityTopic,
		TelemetryInterval: c.telemetryInterval.String(),
		TelemetryTopic:    &c.telemetryTopic,
		TelemetryRetain:   &c.telemetryRetain,
		RealtimeEnable:    &c.realtimeEnable,
		RealtimeTopic:     &c.realtimeTopic,
		RealtimeRetain:    &c.realtimeRetain,
		TopicPrefix:       c.topicPrefix,
		LogDebug:          &c.logDebug,
		LogMessages:       &c.logMessages,
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c ModbusConfig) convertToRead() modbusConfigRead {
	return modbusConfigRead{
		Device:      c.device,
		BaudRate:    c.baudRate,
		ReadTimeout: c.readTimeout.String(),
		LogDebug:    &c.logDebug,
	}
}

func (c DeviceConfig) convertToRead() deviceConfigRead {
	return deviceConfigRead{
		SkipFields:                c.skipFields,
		SkipCategories:            c.skipCategories,
		TelemetryViaMqttClients:   c.telemetryViaMqttClients,
		RealtimeViaMqttClients:    c.realtimeViaMqttClients,
		RestartInterval:           c.restartInterval.String(),
		RestartIntervalMaxBackoff: c.restartIntervalMaxBackoff.String(),
		LogDebug:                  &c.logDebug,
		LogComDebug:               &c.logComDebug,
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c VictronDeviceConfig) convertToRead() victronDeviceConfigRead {
	return victronDeviceConfigRead{
		General: c.DeviceConfig.convertToRead(),
		Device:  c.device,
		Kind:    c.kind.String(),
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c ModbusDeviceConfig) convertToRead() modbusDeviceConfigRead {
	return modbusDeviceConfigRead{
		General: c.DeviceConfig.convertToRead(),
		Bus:     c.bus,
		Kind:    c.kind.String(),
		Address: fmt.Sprintf("0x%02x", c.address),
		Relays: func(inp map[string]RelayConfig) (oup map[string]relayConfigRead) {
			oup = make(map[string]relayConfigRead, len(inp))
			for k, v := range inp {
				oup[k] = v.convertToRead()
			}
			return oup
		}(c.relays),
		PollInterval: c.pollInterval.String(),
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c RelayConfig) convertToRead() relayConfigRead {
	return relayConfigRead{
		Description: &c.description,
		OpenLabel:   &c.openLabel,
		ClosedLabel: &c.closedLabel,
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c HttpDeviceConfig) convertToRead() httpDeviceConfigRead {
	return httpDeviceConfigRead{
		General:      c.DeviceConfig.convertToRead(),
		Url:          c.url.String(),
		Kind:         c.kind.String(),
		Username:     c.username,
		Password:     c.password,
		PollInterval: c.pollInterval.String(),
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c MqttDeviceConfig) convertToRead() mqttDeviceConfigRead {
	return mqttDeviceConfigRead{
		General:     c.DeviceConfig.convertToRead(),
		MqttTopics:  c.mqttTopics,
		MqttClients: c.mqttClients,
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c ViewConfig) convertToRead() viewConfigRead {
	return viewConfigRead{
		Name:         c.name,
		Title:        c.title,
		Devices:      convertListToRead[ViewDeviceConfig, viewDeviceConfigRead](c.devices),
		Autoplay:     &c.autoplay,
		AllowedUsers: maps.Keys(c.allowedUsers),
		Hidden:       &c.hidden,
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c ViewDeviceConfig) convertToRead() viewDeviceConfigRead {
	return viewDeviceConfigRead{
		Name:           c.name,
		Title:          c.title,
		SkipFields:     c.skipFields,
		SkipCategories: c.skipCategories,
	}
}

//lint:ignore U1000 linter does not catch that this is used genric code
func (c HassDiscovery) convertToRead() hassDiscoveryRead {
	return hassDiscoveryRead{
		TopicPrefix:    &c.topicPrefix,
		ViaMqttClients: c.viaMqttClients,
		Devices:        c.devices,
		Categories:     c.categories,
		Registers:      c.registers,
	}
}
