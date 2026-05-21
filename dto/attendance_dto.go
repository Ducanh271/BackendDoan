package dto

type MobileCheckInRequest struct {
	Images    []string `json:"images" binding:"required"`
	BSSID     string   `json:"bssid" binding:"required"`
	Latitude  float64  `json:"latitude" binding:"required"`
	Longitude float64  `json:"longitude" binding:"required"`
}
