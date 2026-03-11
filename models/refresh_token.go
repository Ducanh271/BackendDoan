package models

import "time"

type RefreshToken struct {
	ID         int       `json:"id"`
	EmployeeID int       `json:"employee_id"`
	Token      string    `json:"token"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}
