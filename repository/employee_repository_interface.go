package repository

import (
	"duckanh/backend-doan/models"
	"time"
)

type EmployeeRepository interface {
	Save(employee *models.Employee) error
	GetAll() ([]models.Employee, error)
	GetAllAdmin() ([]models.Employee, error)
	FindByID(ID int) (*models.Employee, error)
	GetByIDAdmin(ID int) (*models.Employee, error)
	LogAccess(log models.AccessLog) error
	FindByEmployeeCode(code string) (*models.Employee, error)
	UpdatePasswordAndStatus(employeeID int, newHash string, isFirstLogin bool) error
	UpdateOTP(employeeID int, otpCode string, expiresAt time.Time) error
	DeleteRefreshToken(employeeID int) error
	GetLastEmployeeCode() (string, error)
	CountByRole(role string) (int, error)
	CreateAdmin(employee *models.Employee) error
}
