package models

import (
	"time"
)

type AccessLog struct {
	ID                 int       `json:"id"`
	EmployeeID         *int      `json:"employee_id"`
	AccessPoint        string    `json:"access_point"`
	ActionTime         time.Time `json:"action_time"`
	Status             string    `json:"status"`
	ConfidenceDistance float64   `json:"confidence_distance"`
	FailureReason      string    `json:"failure_reason"`
}
