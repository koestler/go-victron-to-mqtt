package device

import (
	"github.com/koestler/go-iotdevice/dataflow"
	"strings"
)

type TelemetryMessage struct {
	Time                   string                           `json:"Time"`
	NextTelemetry          string                           `json:"NextTelemetry"`
	Model                  string                           `json:"Model"`
	SecondsSinceLastUpdate float64                          `json:"SecondsSinceLastUpdate"`
	NumericValues          map[string]NumericTelemetryValue `json:"NumericValues,omitempty"`
	TextValues             map[string]TextTelemetryValue    `json:"TextValues,omitempty"`
}

type NumericTelemetryValue struct {
	Category    string  `json:"Cat"`
	Description string  `json:"Desc"`
	Value       float64 `json:"Val"`
	Unit        string  `json:"Unit,omitempty"`
}

type TextTelemetryValue struct {
	Category    string `json:"Cat"`
	Description string `json:"Desc"`
	Value       string `json:"Val"`
}

func convertValuesToNumericTelemetryValues(values []dataflow.Value) (ret map[string]NumericTelemetryValue) {
	ret = make(map[string]NumericTelemetryValue, len(values))

	for _, value := range values {
		if numeric, ok := value.(dataflow.NumericRegisterValue); ok {
			ret[value.Register().Name()] = NumericTelemetryValue{
				Category:    numeric.Register().Category(),
				Description: numeric.Register().Description(),
				Value:       numeric.Value(),
				Unit:        numeric.Register().Unit(),
			}
		}
	}

	return
}

func convertValuesToTextTelemetryValues(values []dataflow.Value) (ret map[string]TextTelemetryValue) {
	ret = make(map[string]TextTelemetryValue, len(values))

	for _, value := range values {
		if text, ok := value.(dataflow.TextRegisterValue); ok {
			ret[value.Register().Name()] = TextTelemetryValue{
				Category:    text.Register().Category(),
				Description: text.Register().Description(),
				Value:       text.Value(),
			}
		}
	}

	return
}

func getTelemetryTopic(topic string, device Device) string {
	// replace Device/Value specific placeholders
	topic = strings.Replace(topic, "%DeviceName%", device.Config().Name(), 1)
	return topic
}
