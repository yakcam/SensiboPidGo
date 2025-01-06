package models

type Configuration struct {
	ApiToken          string
	DeviceId          string
	TargetTemperature float64
	Gain              float64
	Error             int
}
