// service/register_service.go
package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/infrastructure/ai"
	"duckanh/backend-doan/models"
	"duckanh/backend-doan/repository"
)

type RegisterService struct {
	Repo repository.EmployeeRepository
	AI   ai.AIClient
}

// Hàm khởi tạo Service
func NewRegisterService(repo repository.EmployeeRepository, aiClient ai.AIClient) *RegisterService {
	return &RegisterService{
		Repo: repo,
		AI:   aiClient,
	}
}

// Hàm xử lý logic Đăng ký nhân viên
func (s *RegisterService) RegisterEmployee(req dto.RegisterEmployeeRequest) (*dto.RegisterEmployeeResponse, error) {
	// 1. Chuyển đổi mảng ảnh base64 (từ DTO) thành [][]byte để gửi sang AI
	var imageBytes [][]byte
	for _, b64Str := range req.Images {
		// Xử lý trường hợp chuỗi base64 có chứa prefix "data:image/jpeg;base64," (thường gặp nếu gửi từ Web)
		cleanBase64 := b64Str
		if idx := strings.Index(cleanBase64, ","); idx != -1 {
			cleanBase64 = cleanBase64[idx+1:]
		}

		decoded, err := base64.StdEncoding.DecodeString(cleanBase64)
		if err != nil {
			return nil, fmt.Errorf("lỗi giải mã ảnh base64: %v", err)
		}
		imageBytes = append(imageBytes, decoded)
	}

	// 2. Gọi sang AI Service để Check Liveness, Deepfake và trích xuất Embedding
	aiResp, err := s.AI.RegisterFace(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("lỗi kết nối AI: %v", err)
	}

	// 3. Kiểm tra quyết định từ AI (Voting Logic bên Python)
	if aiResp.Status != "success" {
		// Nếu AI chê (Ảnh fake, ảnh in, không đủ mặt,...)
		// Trả về failed cùng với lý do, nhưng không ném ra error hệ thống (để HTTP vẫn trả về status code hợp lý)
		return &dto.RegisterEmployeeResponse{
			Status:  "failed",
			Message: aiResp.Error,
			Stats:   aiResp.Stats,
		}, nil
	}

	// 4. Nếu AI duyệt Pass -> Lắp ráp dữ liệu thành models.Employee
	now := time.Now()
	newEmployee := models.Employee{
		EmployeeCode: req.EmployeeCode,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		DepartmentID: req.DepartmentID,
		PositionID:   req.PositionID,
		Embedding:    aiResp.Embedding, // Vector chuẩn 128 chiều đã được AI tính trung bình
		Status:       "ACTIVE",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 5. Lưu xuống Database
	err = s.Repo.Save(newEmployee)
	if err != nil {
		return nil, fmt.Errorf("lỗi lưu database (có thể trùng mã NV/Email): %v", err)
	}

	// 6. Hoàn thành xuất sắc
	return &dto.RegisterEmployeeResponse{
		Status:  "success",
		Message: "Đăng ký nhân viên thành công!",
		Stats:   aiResp.Stats,
	}, nil
}
