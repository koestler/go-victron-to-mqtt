package main

import (
	"github.com/koestler/go-ve-sensor/vedata"
	"github.com/koestler/go-ve-sensor/vehttp"
	"log"
)

func main() {
	log.Print("start go-ve-sensor...")

	_ = "breakpoint"

	// start http server
	httpdConfig, err := GetHttpdConfig()
	if err == nil {
		log.Print("start http server, config=%v", httpdConfig)
		go func() {
			vehttp.Run(httpdConfig.Bind, httpdConfig.Port, HttpRoutes)
		}()
	} else {
		log.Printf("skip http server, err=%v", err)
	}

	// startup Bmv Device
	log.Print("start devices")
	for _, bmvConfig := range GetBmvConfigs() {
		BmvStart(bmvConfig)
	}

	// run database synchronization routine
	log.Print("start database")
	vedata.Run()

	log.Print("start completed")
	select {}

}
