package dto

type UpdateConfigRequest struct {
	Name      string  `json:"name" binding:"required"`
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	MaxRadius float64 `json:"max_radius" binding:"required"`
}

type AddWifiRequest struct {
	CompanyConfigID int    `json:"company_config_id" binding:"required"`
	BSSID           string `json:"bssid" binding:"required"`
	Description     string `json:"description"`
}
