package models

import "time"

// Define the Location struct
type Location struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	LatLon                []float64 `json:"latLon"`
	Address               []string  `json:"address"`
	Country               *string   `json:"country"`
	CountryAlpha2         string    `json:"countryAlpha2"`
	City                  string    `json:"city"`
	CreateTime            TimeInfo  `json:"createTime"`
	UpdateTime            TimeInfo  `json:"updateTime"`
	Features              []string  `json:"features"`
	GeofenceTriggerRadius int       `json:"geofenceTriggerRadius"`
	Subscription          *string   `json:"subscription"`
	ShareAnalytics        bool      `json:"shareAnalytics"`
	Tariff                *string   `json:"tariff"`
	Currency              string    `json:"currency"`
}

// Define the TimeInfo struct
type TimeInfo struct {
	Time       time.Time `json:"time"`
	SecondsAgo int       `json:"secondsAgo"`
}

// Define the Measurements struct
type Measurements struct {
	Time        TimeInfo `json:"time"`
	Temperature float64  `json:"temperature"`
	Humidity    float64  `json:"humidity"`
	FeelsLike   float64  `json:"feelsLike"`
	Rssi        int      `json:"rssi"`
}

type AcState struct {
	Timestamp         TimeInfo `json:"timestamp"`
	On                bool     `json:"on"`
	Mode              string   `json:"mode"`
	TargetTemperature int      `json:"targetTemperature"`
	TemperatureUnit   string   `json:"temperatureUnit"`
	FanLevel          string   `json:"fanLevel"`
	Swing             string   `json:"swing"`
	HorizontalSwing   string   `json:"horizontalSwing"`
	Light             string   `json:"light"`
}

// Define the Result struct
type Result struct {
	AcState      AcState      `json:"acState"`
	Location     Location     `json:"location"`
	Measurements Measurements `json:"measurements"`
}

// Define the Response struct
type PodsResponse struct {
	Status string `json:"status"`
	Result Result `json:"result"`
}
