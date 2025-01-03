package main

import (
	"SensiboPidGo/apiClient"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"go.einride.tech/pid"
)

func run() int {
	apiToken := os.Getenv("SENSIBO_API_TOKEN")
	if len(apiToken) == 0 {
		log.Fatal("SENSIBO_API_TOKEN is not set")
		return -2
	}

	deviceId := os.Getenv("SENSIBO_DEVICE_ID")
	if len(deviceId) == 0 {
		log.Fatal("SENSIBO_DEVICE_ID is not set")
		return -3
	}

	targetTempString := os.Getenv("TARGET_TEMPERATURE")
	if len(targetTempString) == 0 {
		log.Fatal("TARGET_TEMPERATURE is not set")
		return -4
	}
	targetTemp, _ := strconv.ParseFloat(targetTempString, 32)
	log.Println("Target temperature is:", targetTemp)

	gainString := os.Getenv("GAIN")
	if len(gainString) == 0 {
		log.Fatal("GAIN is not set")
		return -5
	}
	gain, _ := strconv.ParseFloat(gainString, 32)
	log.Println("Gain is:", gain)

	apiResponse, err := apiClient.GetPods(deviceId, apiToken)
	if err != nil {
		log.Fatal(err)
	}

	// Print the latest temperature
	log.Printf("Temperature at %+s was %+v.\n", apiResponse.Result.Measurements.Time.Time, apiResponse.Result.Measurements.Temperature)
	lastResultTime := apiResponse.Result.Measurements.Time.Time
	targetTemperature := targetTemp

	// Create a PID controller.
	c := pid.Controller{
		Config: pid.ControllerConfig{
			ProportionalGain: 7,
			IntegralGain:     0,
			DerivativeGain:   0,
		},
	}

	// Update the PID controller.
	c.Update(pid.ControllerInput{
		ReferenceSignal:  targetTemperature,
		ActualSignal:     apiResponse.Result.Measurements.Temperature,
		SamplingInterval: 0,
	})
	log.Printf("PID State: %+v\n", c.State)

	// Loop round and update
	for {
		apiResponse, err := apiClient.GetPods(deviceId, apiToken)
		if err != nil {
			log.Println(err)
		} else if !apiResponse.Result.AcState.On {
			log.Println("AC is off, waiting 5 minutes before checking again.")
			time.Sleep(300000000000)
		} else if apiResponse.Result.Measurements.Time.Time != lastResultTime {
			c.Update(pid.ControllerInput{
				ReferenceSignal:  targetTemperature,
				ActualSignal:     apiResponse.Result.Measurements.Temperature,
				SamplingInterval: apiResponse.Result.Measurements.Time.Time.Sub(lastResultTime),
			})
			log.Printf("Temperature at %+s was %+v\n", apiResponse.Result.Measurements.Time.Time, apiResponse.Result.Measurements.Temperature)
			log.Printf("PID State: %+v\n", c.State)
			lastResultTime = apiResponse.Result.Measurements.Time.Time
			requestedTemperature := int(math.Round(math.Min(targetTemperature+c.State.ControlSignal, 30.0)))
			if requestedTemperature != apiResponse.Result.AcState.TargetTemperature {
				log.Printf("Setting temperature to %+v\n", requestedTemperature)
				apiClient.SetTemperature(deviceId, apiToken, requestedTemperature)
			} else {
				log.Println("No temperature change needed.")
			}
		} else {
			log.Println("No new data")
		}

		time.Sleep(31000000000)
	}
}

func main() {
	os.Exit(run())
}
