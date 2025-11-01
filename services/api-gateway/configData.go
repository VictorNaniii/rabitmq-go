package main

import "os"

type ConfigPort struct {
	tripPort string
}

func ConfigData() *ConfigPort {
	var configPort ConfigPort

	if configPort.tripPort = os.Getenv("TRIP_SERVICE_PORT"); configPort.tripPort == "" {
		configPort.tripPort = "8083"
	}
	return &ConfigPort{
		tripPort: configPort.tripPort,
	}
}
