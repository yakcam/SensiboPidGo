package main

import (
	"SensiboPidGo/models"
	"log"
	"os"
	"strconv"
)

func GetConfiguration() models.Configuration {
	config := models.Configuration{Error: 0}

	config.ApiToken = os.Getenv("SENSIBO_API_TOKEN")

	config.DeviceId = os.Getenv("SENSIBO_DEVICE_ID")

	targetTempString := os.Getenv("TARGET_TEMPERATURE")
	targetTemp, _ := strconv.ParseFloat(targetTempString, 32)
	config.TargetTemperature = targetTemp

	gainString := os.Getenv("GAIN")
	gain, _ := strconv.ParseFloat(gainString, 32)
	config.Gain = gain

	if len(config.ApiToken) == 0 {
		log.Fatal("SENSIBO_API_TOKEN is not set")
		config.Error = -2
	}

	if len(config.DeviceId) == 0 {
		log.Fatal("SENSIBO_DEVICE_ID is not set")
		config.Error = -3
	}

	if len(targetTempString) == 0 {
		log.Fatal("TARGET_TEMPERATURE is not set")
		config.Error = -4
	}

	if len(gainString) == 0 {
		log.Fatal("GAIN is not set")
		config.Error = -5
	}

	return config
}
