package main

import (
	"SensiboPidGo/apiClient"
	"fmt"
	"math"
	"os"
	"time"

	"go.einride.tech/pid"
)

func run() int {
	apiToken := os.Getenv("SENSIBO_API_TOKEN")
	if len(apiToken) == 0 {
		fmt.Println("SENSIBO_API_TOKEN is not set")
		return -2
	}

	deviceId := os.Getenv("SENSIBO_DEVICE_ID")
	if len(deviceId) == 0 {
		fmt.Println("SENSIBO_DEVICE_ID is not set")
		return -3
	}

	apiResponse, err := apiClient.GetPods(deviceId, apiToken)
	if err != nil {
		fmt.Println(err)
	}

	// Print the latest temperature
	fmt.Printf("%+s: %+v\n", apiResponse.Result.Measurements.Time.Time, apiResponse.Result.Measurements.Temperature)
	lastResultTime := apiResponse.Result.Measurements.Time.Time
	targetTemperature := 21.0
	requestedTemperature := int(targetTemperature)

	// Create a PID controller.
	c := pid.Controller{
		Config: pid.ControllerConfig{
			ProportionalGain: 3.0,
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
	fmt.Printf("%+v\n", c.State)

	// Loop round and update
	for {
		apiResponse, err := apiClient.GetPods(deviceId, apiToken)
		if err != nil {
			fmt.Println(err)
		} else if apiResponse.Result.Measurements.Time.Time != lastResultTime {
			c.Update(pid.ControllerInput{
				ReferenceSignal:  targetTemperature,
				ActualSignal:     apiResponse.Result.Measurements.Temperature,
				SamplingInterval: apiResponse.Result.Measurements.Time.Time.Sub(lastResultTime),
			})
			fmt.Printf("%+s: %+v\n", apiResponse.Result.Measurements.Time.Time, apiResponse.Result.Measurements.Temperature)
			fmt.Printf("%+v\n", c.State)
			lastResultTime = apiResponse.Result.Measurements.Time.Time
			requestedTemperature = int(math.Round(math.Min(targetTemperature+c.State.ControlSignal, 30.0)))
			if requestedTemperature != int(apiResponse.Result.Measurements.Temperature) {
				fmt.Printf("Setting temperature to %+v\n", requestedTemperature)
				apiClient.SetTemperature(deviceId, apiToken, requestedTemperature)
			}
		} else {
			fmt.Println("No new data")
		}

		time.Sleep(29000000000)
	}
}

func main() {
	os.Exit(run())
}
