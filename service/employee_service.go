package service

import (
	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/repository"
	"time"
)

type EmployeeService struct {
	EmpRepo repository.EmployeeRepository
}

func NewEmployeeService(empRepo repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		EmpRepo: empRepo,
	}
}

func (s *EmployeeService) GetProfileByID(id int) (*dto.ProfileResponse, error) {
	emp, err := s.EmpRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Xử lý an toàn các trường con trỏ (Tránh lỗi nil pointer dereference)
	email := ""
	if emp.Email != nil {
		email = *emp.Email
	}

	phone := ""
	if emp.Phone != nil {
		phone = *emp.Phone
	}

	deptID := 0
	if emp.DepartmentID != nil {
		deptID = *emp.DepartmentID
	}

	posID := 0
	if emp.PositionID != nil {
		posID = *emp.PositionID
	}

	hireDateStr := ""
	if emp.HireDate != nil {
		hireDateStr = emp.HireDate.Format("2006-01-02")
	}

	hireDate, err := time.Parse("2006-01-02", hireDateStr)
	return &dto.ProfileResponse{
		ID:           emp.ID,
		EmployeeCode: emp.EmployeeCode,
		Name:         emp.Name,
		Email:        email,
		Phone:        phone,
		DepartmentID: deptID,
		PositionID:   posID,
		HireDate:     &hireDate,
		Status:       emp.Status,
	}, nil
}
