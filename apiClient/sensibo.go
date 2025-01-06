package apiClient

import (
	"SensiboPidGo/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SetTemperature(deviceId string, apiToken string, temperature int) error {
	return SetProperty(deviceId, apiToken, "targetTemperature", temperature)
}

func SetMode(deviceId string, apiToken string, mode string) error {
	return SetProperty(deviceId, apiToken, "mode", mode)
}

func SetProperty[T int | string](deviceId string, apiToken string, propertyName string, propertyValue T) error {

	// Create the JSON payload
	var payload string
	switch any(propertyValue).(type) {
	case string:
		payload = fmt.Sprintf(`{"newValue": "%s"}`, propertyValue)
	case int:
		payload = fmt.Sprintf(`{"newValue": %d}`, propertyValue)
	default:
		return errors.New("unsupported type")
	}

	deviceUrl := fmt.Sprintf("https://home.sensibo.com/api/v2/pods/%s/acStates/%s?apiKey=%s", deviceId, propertyName, apiToken)

	// Create a new request
	req, err := http.NewRequest(http.MethodPatch, deviceUrl, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Add the JSON payload to the request
	req.Header.Add("Content-Type", "application/json")
	req.Body = io.NopCloser(strings.NewReader(payload))

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return err
	}

	return nil
}

func GetPods(deviceId string, apiToken string) (models.PodsResponse, error) {
	deviceUrl := fmt.Sprintf("https://home.sensibo.com/api/v2/pods/%s?apiKey=%s&fields=location,measurements,acState", deviceId, apiToken)

	resp, err := http.Get(deviceUrl)
	if err != nil {
		fmt.Println("Error making request:", err)
		return models.PodsResponse{}, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return models.PodsResponse{}, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return models.PodsResponse{}, err
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResponse models.PodsResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return models.PodsResponse{}, err
	}

	return apiResponse, nil
}
