package models

type CompanyWifi struct {
	ID              int    `json:"id"`
	CompanyConfigID int    `json:"company_config_id"`
	BSSID           string `json:"bssid"`
	Description     string `json:"description"`
}
