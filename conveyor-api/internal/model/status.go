package model

import "time"

type SystemStatus struct {
	LastSeen         time.Time `json:"last_seen"`
	BeltRunning      bool      `json:"belt_running"`
	FanOn            bool      `json:"fan_on"`
	BuzzerOn         bool      `json:"buzzer_on"`
	DoorAngle        int       `json:"door_angle"`
	ObjectCountTotal int64     `json:"object_count_total"`
	ActiveAlerts     []Alert   `json:"active_alerts"`
}
