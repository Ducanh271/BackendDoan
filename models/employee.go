package models

import "time"

type Employee struct {
	ID           int        `json:"id"`
	EmployeeCode string     `json:"employee_code"`
	Name         string     `json:"name"`
	Email        *string    `json:"email"`
	Phone        *string    `json:"phone"` // Có thể NULL
	DepartmentID *int       `json:"department_id"`
	PositionID   *int       `json:"position_id"`
	HireDate     *time.Time `json:"hire_date"`
	Embedding    []float64  `json:"-"` // Cột này chỉ dùng ở Backend, không trả JSON về Frontend
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
