package model

import "time"

type AlertType string

const (
	AlertGas      AlertType = "GAS"
	AlertHumidity AlertType = "HUMIDITY"
)

type Alert struct {
	ID           string     `json:"id"`
	Timestamp    time.Time  `json:"timestamp"`
	Type         AlertType  `json:"type"`
	TriggerValue int        `json:"trigger_value"`
	Threshold    int        `json:"threshold"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
	Active       bool       `json:"active"`
}
