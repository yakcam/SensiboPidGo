package apiClient

import (
	"SensiboPidGo/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetPods(deviceId string, apiToken string) (models.PodsResponse, error) {
	deviceUrl := fmt.Sprintf("https://home.sensibo.com/api/v2/pods/%s?apiKey=%s&fields=location,measurements", deviceId, apiToken)

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
