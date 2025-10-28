package main

import (
	"SensiboPidGo/apiClient"
	"log"
	"math"
	"os"
	"time"

	"go.einride.tech/pid"
)

func run() int {
	config := GetConfiguration()

	if config.Error != 0 {
		return config.Error
	}

	log.Println("Target temperature is:", config.TargetTemperature)
	log.Println("Gain is:", config.Gain)

	apiResponse, err := apiClient.GetPods(config.DeviceId, config.ApiToken)
	if err != nil {
		log.Fatal(err)
	}

	// Print the latest temperature
	log.Printf("Temperature at %+s was %+v.\n", apiResponse.Result.Measurements.Time.Time, apiResponse.Result.Measurements.Temperature)
	lastResultTime := apiResponse.Result.Measurements.Time.Time

	// Create a PID controller.
	c := pid.Controller{
		Config: pid.ControllerConfig{
			ProportionalGain: config.Gain,
			IntegralGain:     0,
			DerivativeGain:   0,
		},
	}

	// Update the PID controller.
	c.Update(pid.ControllerInput{
		ReferenceSignal:  config.TargetTemperature,
		ActualSignal:     apiResponse.Result.Measurements.Temperature,
		SamplingInterval: 0,
	})
	log.Printf("PID State: %+v\n", c.State)

	// Loop round and update
	for {
		apiResponse, err := apiClient.GetPods(config.DeviceId, config.ApiToken)
		if err != nil {
			log.Println(err)
		} else if !apiResponse.Result.AcState.On {
			log.Println("AC is off, waiting 5 minutes before checking again.")
			time.Sleep(300000000000)
		} else if apiResponse.Result.AcState.Mode == "cool" {
			log.Println("AC is cooling, waiting 5 minutes before checking again.")
			time.Sleep(300000000000)
		} else if apiResponse.Result.Measurements.Time.Time != lastResultTime {
			c.Update(pid.ControllerInput{
				ReferenceSignal:  config.TargetTemperature,
				ActualSignal:     apiResponse.Result.Measurements.Temperature,
				SamplingInterval: apiResponse.Result.Measurements.Time.Time.Sub(lastResultTime),
			})
			log.Printf("Temperature at %+s was %+v\n", apiResponse.Result.Measurements.Time.Time, apiResponse.Result.Measurements.Temperature)
			log.Printf("PID State: %+v\n", c.State)
			lastResultTime = apiResponse.Result.Measurements.Time.Time

			// Temperature control
			var requestedTemperature int
			if c.State.ControlSignal >= 0 {
				requestedTemperature = int(math.Round(math.Min(config.TargetTemperature+c.State.ControlSignal, 30.0))) // Max ac temp
			} else {
				// When temperature is reached or exceeded, set to minimum temperature (17°C)
				requestedTemperature = 17
			}
			if requestedTemperature != apiResponse.Result.AcState.TargetTemperature {
				log.Printf("Setting temperature to %+v\n", requestedTemperature)
				apiClient.SetTemperature(config.DeviceId, config.ApiToken, requestedTemperature)
			} else {
				log.Println("No temperature change needed.")
			}

			// Fan control based on whether heat is required
			var requestedFanLevel string
			if c.State.ControlSignal >= 0 {
				requestedFanLevel = "high" // Heat is required
			} else {
				requestedFanLevel = "auto" // Temperature reached
			}

			if apiResponse.Result.AcState.FanLevel != requestedFanLevel {
				log.Printf("Setting fan level to %s\n", requestedFanLevel)
				apiClient.SetFanLevel(config.DeviceId, config.ApiToken, requestedFanLevel)
			} else {
				log.Println("No fan level change needed.")
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
