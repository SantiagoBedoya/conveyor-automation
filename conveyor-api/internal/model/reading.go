package model

import "time"

type SensorReading struct {
	ID            string    `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	GasValue      int       `json:"gas_value"`
	HumidityValue int       `json:"humidity_value"`
	DistanceCm    float64   `json:"distance_cm"`
	ObjectCount   int       `json:"object_count"`
	BeltRunning   bool      `json:"belt_running"`
	FanOn         bool      `json:"fan_on"`
	BuzzerOn      bool      `json:"buzzer_on"`
	DoorAngle     int       `json:"door_angle"`
}
