package victron

import (
	"github.com/koestler/go-iotdevice/dataflow"
)

var RegisterListBmvProduct = dataflow.Registers{
	dataflow.CreateNumberRegisterStruct(
		"Product",
		"ProductId",
		"Product Id",
		0x0100,
		false,
		1,
		"",
	),
	dataflow.CreateNumberRegisterStruct(
		"Product",
		"ProductRevision",
		"Product revision",
		0x0101,
		false,
		1,
		"",
	),
	dataflow.CreateTextRegisterStruct(
		"Product",
		"SerialNumber",
		"Serial number",
		0x010A,
	),
	dataflow.CreateTextRegisterStruct(
		"Product",
		"ModelName",
		"Model name",
		0x010B,
	),
	dataflow.CreateTextRegisterStruct(
		"Product",
		"Description",
		"Description",
		0x010C,
	),
	dataflow.CreateNumberRegisterStruct(
		"Product",
		"Uptime",
		"Device uptime",
		0x0120,
		false,
		1,
		"s",
	),
	// skipped Bluetooth capabilities
}

var RegisterListBmvMonitor = dataflow.Registers{
	dataflow.CreateNumberRegisterStruct(
		"Essential",
		"MainVoltage",
		"Main Voltage",
		0xED8D,
		true,
		0.01,
		"V",
	),
	dataflow.CreateNumberRegisterStruct(
		"Monitor",
		"AuxVoltage",
		"Aux (starter) Voltage",
		0xED7D,
		false,
		0.01,
		"V",
	),
	dataflow.CreateNumberRegisterStruct(
		"Essential",
		"Current",
		"Current",
		0xED8F,
		true,
		0.1,
		"A",
	),
	dataflow.CreateNumberRegisterStruct(
		"Monitor",
		"CurrentHighRes",
		"Current (high resolution)",
		0xED8C,
		true,
		0.001,
		"A",
	),
	dataflow.CreateNumberRegisterStruct(
		"Essential",
		"Power",
		"Power",
		0xED8E,
		true,
		1,
		"W",
	), dataflow.CreateNumberRegisterStruct(
		"Monitor",
		"Consumed",
		"Consumed Ah",
		0xEEFF,
		true,
		0.1,
		"Ah",
	), dataflow.CreateNumberRegisterStruct(
		"Essential",
		"SOC",
		"State Of Charge",
		0x0FFF,
		false,
		0.01,
		"%",
	), dataflow.CreateNumberRegisterStruct(
		"Monitor",
		"TTG",
		"Time to go",
		0x0FFE,
		false,
		1,
		"min",
	), dataflow.CreateNumberRegisterStruct(
		"Monitor",
		"Temperature",
		"Temperature",
		0xEDEC,
		false,
		0.01,
		"K",
	), dataflow.CreateNumberRegisterStruct(
		"Monitor",
		"MidPointVoltage",
		"Mid-point voltage",
		0x0382,
		false,
		0.01,
		"V",
	), dataflow.CreateNumberRegisterStruct(
		"Monitor",
		"MidPointVoltageDeviation",
		"Mid-point voltage deviation",
		0x0383,
		true,
		0.1,
		"%",
	), dataflow.CreateNumberRegisterStruct(
		"Monitor",
		"SynchronizationState",
		"Synchronization state",
		0xEEB6,
		false,
		1,
		"1",
		// todo: decode this enum
	),
}

var RegisterListBmvHistoric = dataflow.Registers{
	dataflow.CreateNumberRegisterStruct(
		"Historic",
		"DepthOfTheDeepestDischarge",
		"Depth of the deepest discharge",
		0x0300,
		true,
		0.1,
		"Ah",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"DepthOfTheLastDischarge",
		"Depth of the last discharge",
		0x0301,
		true,
		0.1,
		"Ah",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"DepthOfTheAverageDischarge",
		"Depth of the average discharge",
		0x0302,
		true,
		0.1,
		"Ah",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"NumberOfCycles",
		"Number of cycles",
		0x0303,
		false,
		1,
		"",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"NumberOfFullDischarges",
		"Number of full discharges",
		0x0304,
		false,
		1,
		"",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"CumulativeAmpHours",
		"Cumulative Amp Hours",
		0x0305,
		true,
		0.1,
		"Ah",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"MainVoltageMinimum",
		"Minimum Voltage",
		0x0306,
		false,
		0.01,
		"V",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"MainVoltageMaximum",
		"Maximum Voltage",
		0x0307,
		false,
		0.01,
		"V",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"HoursSinceFullCharge",
		"Hours since full charge",
		0x0308,
		false,
		float64(24)/float64(86400),
		"h",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"NumberOfAutomaticSynchronizations",
		"Number of automatic synchronizations",
		0x0309,
		false,
		1,
		"",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"NumberOfLowMainVoltageAlarms",
		"Number of Low Voltage Alarms",
		0x030A,
		false,
		1,
		"",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"NumberOfHighMainVoltageAlarms",
		"Number of High Voltage Alarms",
		0x030B,
		false,
		1,
		"",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"AuxVoltageMinimum",
		"Minimum Starter Voltage",
		0x030E,
		true,
		0.01,
		"V",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"AuxVoltageMaximum",
		"Maximum Starter Voltage",
		0x030F,
		true,
		0.01,
		"V",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"AmountOfDischargedEnergy",
		"Amount of discharged energy",
		0x0310,
		false,
		0.01,
		"kWh",
	), dataflow.CreateNumberRegisterStruct(
		"Historic",
		"AmountOfChargedEnergy",
		"Amount of charged energy",
		0x0311,
		false,
		0.01,
		"kWh",
	),
}

var RegisterListBmv712 = dataflow.MergeRegisters(
	RegisterListBmvProduct,
	RegisterListBmvMonitor,
	RegisterListBmvHistoric,
)

var RegisterListBmv702 = dataflow.FilterRegisters(
	RegisterListBmv712,
	[]string{
		"ProductRevision",
		"Description",
	},
)

var RegisterListBmv700 = dataflow.FilterRegisters(
	RegisterListBmv702,
	[]string{
		"AuxVoltage",
		"Temperature",
		"MidPointVoltage",
		"MidPointVoltageDeviation",
		"AuxVoltageMinimum",
		"AuxVoltageMaximum",
	},
)
