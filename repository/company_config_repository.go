package repository

import (
	"database/sql"
	"duckanh/backend-doan/models"
	"fmt"
)

type CompanyConfigRepository interface {
	GetActiveConfig() (*models.CompanyConfig, error)
	UpdateConfig(cfg *models.CompanyConfig) error // MỚI
	AddWifi(wifi *models.CompanyWifi) error       // MỚI
	DeleteWifi(wifiID int) error                  // MỚI
}

type MySQLCompanyConfigRepository struct {
	DB *sql.DB
}

func NewMySQLCompanyConfigRepository(db *sql.DB) CompanyConfigRepository {
	return &MySQLCompanyConfigRepository{DB: db}
}

func (r *MySQLCompanyConfigRepository) GetActiveConfig() (*models.CompanyConfig, error) {
	// 1. Lấy thông tin tọa độ công ty (Giả sử hệ thống hiện tại chỉ có 1 trụ sở ID = 1)
	queryConfig := `SELECT id, name, latitude, longitude, max_radius FROM company_configs ORDER BY id ASC LIMIT 1`
	var cfg models.CompanyConfig

	err := r.DB.QueryRow(queryConfig).Scan(
		&cfg.ID,
		&cfg.Name,
		&cfg.Latitude,
		&cfg.Longitude,
		&cfg.MaxRadius,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("chưa có cấu hình công ty nào trong hệ thống")
		}
		return nil, fmt.Errorf("lỗi lấy cấu hình công ty: %v", err)
	}

	// 2. Lấy danh sách Wi-Fi thuộc về công ty này
	queryWifis := `SELECT id, company_config_id, bssid, description FROM company_wifis WHERE company_config_id = ?`
	rows, err := r.DB.Query(queryWifis, cfg.ID)
	if err != nil {
		return nil, fmt.Errorf("lỗi lấy danh sách Wi-Fi: %v", err)
	}
	defer rows.Close()

	var wifis []models.CompanyWifi
	for rows.Next() {
		var wifi models.CompanyWifi
		var desc sql.NullString

		if err := rows.Scan(&wifi.ID, &wifi.CompanyConfigID, &wifi.BSSID, &desc); err == nil {
			if desc.Valid {
				wifi.Description = desc.String
			}
			wifis = append(wifis, wifi)
		}
	}

	// 3. Gắn mảng Wi-Fi vào struct Config
	cfg.Wifis = wifis
	return &cfg, nil
}
func (r *MySQLCompanyConfigRepository) UpdateConfig(cfg *models.CompanyConfig) error {
	query := `UPDATE company_configs SET name = ?, latitude = ?, longitude = ?, max_radius = ? WHERE id = ?`
	_, err := r.DB.Exec(query, cfg.Name, cfg.Latitude, cfg.Longitude, cfg.MaxRadius, cfg.ID)
	return err
}

// 2. Thêm một mạng Wi-Fi mới
func (r *MySQLCompanyConfigRepository) AddWifi(wifi *models.CompanyWifi) error {
	query := `INSERT INTO company_wifis (company_config_id, bssid, description) VALUES (?, ?, ?)`
	result, err := r.DB.Exec(query, wifi.CompanyConfigID, wifi.BSSID, wifi.Description)
	if err != nil {
		return fmt.Errorf("lỗi thêm Wi-Fi (có thể BSSID bị trùng): %v", err)
	}

	id, _ := result.LastInsertId()
	wifi.ID = int(id)
	return nil
}

// 3. Xóa một mạng Wi-Fi
func (r *MySQLCompanyConfigRepository) DeleteWifi(wifiID int) error {
	query := `DELETE FROM company_wifis WHERE id = ?`
	_, err := r.DB.Exec(query, wifiID)
	return err
}
