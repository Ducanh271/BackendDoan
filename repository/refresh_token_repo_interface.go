package repository

import "duckanh/backend-doan/models"

type RefreshTokenRepository interface {
	Save(token *models.RefreshToken) error
	FindByToken(tokenString string) (*models.RefreshToken, error)
	DeleteByToken(tokenString string) error
	DeleteAllByEmployeeID(employeeID int) error // Dùng khi muốn "Đăng xuất khỏi mọi thiết bị"
}
