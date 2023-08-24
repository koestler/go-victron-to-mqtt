package dataflow_test

import (
	"fmt"
	"github.com/koestler/go-iotdevice/dataflow"
)

func getSimpleTestRegister(category, name string) dataflow.RegisterStruct {
	return dataflow.NewRegisterStruct(
		category,
		name,
		"",
		dataflow.NumberRegister,
		map[int]string{},
		"",
		40,
		false,
	)
}
func fillSetA(storage *dataflow.ValueStorageInstance) {
	storage.Fill(dataflow.NewNumericRegisterValue(
		"device-0",
		getSimpleTestRegister("set-a", "register-a"),
		0,
	))

	storage.Fill(dataflow.NewNumericRegisterValue(
		"device-0",
		getSimpleTestRegister("set-a", "register-a"),
		1,
	))

	// filling the same register multiple times must not make a difference
	for i := 0; i < 10; i += 1 {
		storage.Fill(dataflow.NewNumericRegisterValue(
			"device-0",
			getSimpleTestRegister("set-a", "register-b"),
			10,
		))

		storage.Fill(dataflow.NewNumericRegisterValue(
			"device-1",
			getSimpleTestRegister("set-a", "register-a"),
			100,
		))
	}
}

func fillSetB(storage *dataflow.ValueStorageInstance) {
	storage.Fill(dataflow.NewNumericRegisterValue(
		"device-1",
		getSimpleTestRegister("set-b", "register-a"),
		101,
	))

	storage.Fill(dataflow.NewNumericRegisterValue(
		"device-2",
		getSimpleTestRegister("set-b", "register-a"),
		200,
	))
}

func fillSetC(storage *dataflow.ValueStorageInstance) {
	for i := 0; i < 1000; i += 1 {
		storage.Fill(dataflow.NewNumericRegisterValue(
			"device-3",
			getSimpleTestRegister("set-c", fmt.Sprintf("register-%d", i)),
			float64(i),
		))
	}
}
