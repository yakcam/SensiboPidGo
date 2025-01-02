package main

import (
	"SensiboPidGo/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

	deviceUrl := fmt.Sprintf("https://home.sensibo.com/api/v2/pods/%s?apiKey=%s&fields=location,measurements", deviceId, apiToken)

	resp, err := http.Get(deviceUrl)
	if err != nil {
		fmt.Println("Error making request:", err)
		return -5
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return -10
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return -20
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResponse models.PodsResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return -30
	}

	// Print the latest temperature
	fmt.Printf("%+s: %+v\n", apiResponse.Result.Measurements.Time.Time, apiResponse.Result.Measurements.Temperature)

	return 0
}

func main() {
	os.Exit(run())
}
