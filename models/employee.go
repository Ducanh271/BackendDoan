package models

import "time"

type Employee struct {
	ID           int        `json:"id"`
	EmployeeCode string     `json:"employee_code"`
	Name         string     `json:"name"`
	Email        *string    `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	IsFirstLogin bool       `json:"is_first_login"`
	OTPCode      *string    `json:"otp_code"`
	OTPExpiresAt *time.Time `json:"otp_expires_at"`
	Phone        *string    `json:"phone"`
	DepartmentID *int       `json:"department_id"`
	PositionID   *int       `json:"position_id"`
	HireDate     *time.Time `json:"hire_date"`
	Embedding    []float64  `json:"-"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
