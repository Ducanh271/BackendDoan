package models

import "time"

type CompanyConfig struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	MaxRadius float64       `json:"max_radius"`
	Wifis     []CompanyWifi `json:"wifis"` // Chứa danh sách các Wi-Fi thuộc công ty này
	UpdatedAt time.Time     `json:"updated_at"`
}
