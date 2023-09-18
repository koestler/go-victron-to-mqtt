package config

import (
	"net/url"
	"regexp"
	"time"
)

type Nameable interface {
	Name() string
}

type Config struct {
	version         int
	projectTitle    string
	logConfig       bool
	logWorkerStart  bool
	logStorageDebug bool
	httpServer      HttpServerConfig
	authentication  AuthenticationConfig
	mqttClients     []MqttClientConfig
	modbus          []ModbusConfig
	devices         []DeviceConfig
	victronDevices  []VictronDeviceConfig
	modbusDevices   []ModbusDeviceConfig
	httpDevices     []HttpDeviceConfig
	mqttDevices     []MqttDeviceConfig
	views           []ViewConfig
	hassDiscovery   []HassDiscovery
}

type HttpServerConfig struct {
	enabled         bool
	bind            string
	port            int
	logRequests     bool
	frontendProxy   *url.URL
	frontendPath    string
	frontendExpires time.Duration
	configExpires   time.Duration
	logDebug        bool
}

type AuthenticationConfig struct {
	enabled           bool
	jwtSecret         []byte
	jwtValidityPeriod time.Duration
	htaccessFile      string
}

type MqttClientConfig struct {
	name            string
	broker          *url.URL
	protocolVersion int

	user     string
	password string
	clientId string

	qos               byte
	keepAlive         time.Duration
	connectRetryDelay time.Duration
	connectTimeout    time.Duration
	topicPrefix       string
	readOnly          bool
	maxBacklogSize    int

	availabilityEnable bool
	availabilityTopic  string
	availabilityRetain bool

	structureEnable   bool
	structureTopic    string
	structureInterval time.Duration
	structureRetain   bool

	telemetryEnable   bool
	telemetryTopic    string
	telemetryInterval time.Duration
	telemetryRetain   bool

	realtimeEnable   bool
	realtimeTopic    string
	realtimeInterval time.Duration
	realtimeRepeat   bool
	realtimeRetain   bool

	logDebug    bool
	logMessages bool
}

type ModbusConfig struct {
	name        string
	device      string
	baudRate    int
	readTimeout time.Duration
	logDebug    bool
}

type DeviceConfig struct {
	name                      string
	skipFields                []string
	skipCategories            []string
	telemetryViaMqttClients   []string
	realtimeViaMqttClients    []string
	restartInterval           time.Duration
	restartIntervalMaxBackoff time.Duration
	logDebug                  bool
	logComDebug               bool
}

type VictronDeviceConfig struct {
	DeviceConfig
	device string
	kind   VictronDeviceKind
}

type ModbusDeviceConfig struct {
	DeviceConfig
	bus          string
	kind         ModbusDeviceKind
	address      byte
	relays       map[string]RelayConfig
	pollInterval time.Duration
}

type RelayConfig struct {
	description string
	openLabel   string
	closedLabel string
}

type HttpDeviceConfig struct {
	DeviceConfig
	url          *url.URL
	kind         HttpDeviceKind
	username     string
	password     string
	pollInterval time.Duration
}

type MqttDeviceConfig struct {
	DeviceConfig
	mqttTopics  []string
	mqttClients []string
}

type ViewConfig struct {
	name         string
	title        string
	devices      []ViewDeviceConfig
	autoplay     bool
	allowedUsers map[string]struct{}
	hidden       bool
}

type ViewDeviceConfig struct {
	name           string
	title          string
	skipFields     []string
	skipCategories []string
}

type HassDiscovery struct {
	topicPrefix       string
	viaMqttClients    []string
	devices           []string
	categories        []string
	categoriesMatcher []*regexp.Regexp
	registers         []string
	registersMatcher  []*regexp.Regexp
}

type VictronDeviceKind int

const (
	VictronUndefinedKind VictronDeviceKind = iota
	VictronRandomBmvKind
	VictronRandomSolarKind
	VictronVedirectKind
)

type ModbusDeviceKind int

const (
	ModbusUndefinedKind ModbusDeviceKind = iota
	ModbusWaveshareRtuRelay8Kind
)

type HttpDeviceKind int

const (
	HttpUndefinedKind HttpDeviceKind = iota
	HttpTeracomKind
	HttpShellyEm3Kind
)
