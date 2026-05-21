package service

import (
	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/models"
	"duckanh/backend-doan/repository"
)

type ConfigService struct {
	ConfigRepo repository.CompanyConfigRepository
}

func NewConfigService(repo repository.CompanyConfigRepository) *ConfigService {
	return &ConfigService{ConfigRepo: repo}
}

// Lấy cấu hình
func (s *ConfigService) GetConfig() (*models.CompanyConfig, error) {
	return s.ConfigRepo.GetActiveConfig()
}

// Cập nhật cấu hình
func (s *ConfigService) UpdateConfig(id int, req dto.UpdateConfigRequest) error {
	cfg := models.CompanyConfig{
		ID:        id,
		Name:      req.Name,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		MaxRadius: req.MaxRadius,
	}
	return s.ConfigRepo.UpdateConfig(&cfg)
}

// Thêm Wi-Fi
func (s *ConfigService) AddWifi(req dto.AddWifiRequest) error {
	wifi := models.CompanyWifi{
		CompanyConfigID: req.CompanyConfigID,
		BSSID:           req.BSSID,
		Description:     req.Description,
	}
	return s.ConfigRepo.AddWifi(&wifi)
}

// Xóa Wi-Fi
func (s *ConfigService) DeleteWifi(id int) error {
	return s.ConfigRepo.DeleteWifi(id)
}
