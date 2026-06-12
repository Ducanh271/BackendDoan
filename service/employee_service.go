package service

import (
	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/models"
	"duckanh/backend-doan/repository"
	"duckanh/backend-doan/security"
	"fmt"
	"log"
	"time"
)

// Thông tin tài khoản admin mặc định được seed khi hệ thống chưa có admin nào.
const (
	DefaultAdminCode     = "admin"
	DefaultAdminName     = "Administrator"
	DefaultAdminPassword = "123456a@"
)

type EmployeeService struct {
	EmpRepo repository.EmployeeRepository
}

func NewEmployeeService(empRepo repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		EmpRepo: empRepo,
	}
}

// EnsureDefaultAdmin tạo tài khoản admin mặc định (admin / 123456a@) nếu hệ thống
// chưa có bất kỳ tài khoản ADMIN nào. Gọi lúc khởi động server.
func (s *EmployeeService) EnsureDefaultAdmin() error {
	count, err := s.EmpRepo.CountByRole("ADMIN")
	if err != nil {
		return fmt.Errorf("không kiểm tra được tài khoản admin: %w", err)
	}
	if count > 0 {
		return nil // Đã có admin, không cần seed
	}

	hash, err := security.HashPassword(DefaultAdminPassword)
	if err != nil {
		return fmt.Errorf("không hash được mật khẩu admin mặc định: %w", err)
	}

	admin := &models.Employee{
		EmployeeCode: DefaultAdminCode,
		Name:         DefaultAdminName,
		PasswordHash: hash,
		Role:         "ADMIN",
		IsFirstLogin: false, // Đăng nhập được ngay, không cần đổi mật khẩu qua OTP
		Status:       "ACTIVE",
	}
	if err := s.EmpRepo.CreateAdmin(admin); err != nil {
		return fmt.Errorf("không tạo được tài khoản admin mặc định: %w", err)
	}

	log.Printf("Đã tạo tài khoản admin mặc định: %s / %s (vui lòng đổi mật khẩu sau khi đăng nhập)", DefaultAdminCode, DefaultAdminPassword)
	return nil
}

func (s *EmployeeService) GetAllEmployees() ([]dto.ProfileResponse, error) {
	employees, err := s.EmpRepo.GetAllAdmin()
	if err != nil {
		return nil, err
	}

	var res []dto.ProfileResponse
	for _, emp := range employees {
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
		hireDate := time.Time{}
		if emp.HireDate != nil {
			hireDate = *emp.HireDate
		}

		res = append(res, dto.ProfileResponse{
			ID:           emp.ID,
			EmployeeCode: emp.EmployeeCode,
			Name:         emp.Name,
			Email:        email,
			Phone:        phone,
			DepartmentID: deptID,
			PositionID:   posID,
			HireDate:     &hireDate,
			Status:       emp.Status,
		})
	}
	return res, nil
}

func (s *EmployeeService) GetEmployeeDetail(id int) (*dto.ProfileResponse, error) {
	emp, err := s.EmpRepo.GetByIDAdmin(id)
	if err != nil {
		return nil, err
	}

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
	hireDate := time.Time{}
	if emp.HireDate != nil {
		hireDate = *emp.HireDate
	}

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

func (s *EmployeeService) GetProfileByID(id int) (*dto.ProfileResponse, error) {
	emp, err := s.EmpRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

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
	hireDate := time.Time{}
	if emp.HireDate != nil {
		hireDate = *emp.HireDate
	}

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




