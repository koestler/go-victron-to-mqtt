package httpServer

var wsRoutes = WsRoutes{
	WsRoute{
		"ws-test",
		"/ws/v0/RoundedValues",
		HandleWsRoundedValues,
	},
}

var httpRoutes = HttpRoutes{
	HttpRoute{
		"DeviceIndex",
		"GET",
		"/api/v0/Devices",
		HandleDeviceIndex,
	},
	HttpRoute{
		"FrontendConfig",
		"GET",
		"/api/v0/FrontendConfig",
		HandleFrontendConfig,
	},
	HttpRoute{
		"RoundedValues",
		"GET",
		"/api/v0/Device/{DeviceId:[a-zA-Z0-9\\-]{1,32}}/RoundedValues",
		HandleDeviceGetRoundedValues,
	},
	HttpRoute{
		"HassMqttSensorsYaml",
		"GET",
		"/api/v0/Hass/MqttSensors",
		HandleHassMqttSensorsYaml,
	},
	HttpRoute{
		"DeviceRoundedValuesWebSocket",
		"GET",
		"/api/v0/ws/RoundedValues",
		HandleWsRoundedValues,
	},
	HttpRoute{
		"ApiIndex",
		"GET",
		"/api{Path:.*}",
		HandleApiNotFound,
	},
}
