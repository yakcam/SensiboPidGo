package main

import (
	"bufio"
	"fmt"
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
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return 0
}

func main() {
	os.Exit(run())
}
