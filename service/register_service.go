// service/register_service.go
package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/infrastructure/ai"
	"duckanh/backend-doan/models"
	"duckanh/backend-doan/repository"
	"duckanh/backend-doan/security"
	"duckanh/backend-doan/utils"
)

type RegisterService struct {
	Repo         repository.EmployeeRepository
	AI           ai.AIClient
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string
}

// Hàm khởi tạo Service
func NewRegisterService(repo repository.EmployeeRepository, aiClient ai.AIClient, smtpHost, smtpPort, smtpEmail, smtpPassword string) *RegisterService {
	return &RegisterService{
		Repo:         repo,
		AI:           aiClient,
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPEmail:    smtpEmail,
		SMTPPassword: smtpPassword,
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

	// 4. Nếu AI duyệt Pass -> Tự động tạo Mã nhân viên và Mật khẩu
	employeeCode, err := s.generateEmployeeCode()
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo mã nhân viên: %v", err)
	}

	rawPassword := s.generateRandomPassword(8)
	hashedPassword, err := security.HashPassword(rawPassword)
	if err != nil {
		return nil, fmt.Errorf("lỗi mã hóa mật khẩu: %v", err)
	}

	now := time.Now()
	newEmployee := models.Employee{
		EmployeeCode: employeeCode,
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		DepartmentID: req.DepartmentID,
		PositionID:   req.PositionID,
		Embedding:    aiResp.Embedding, // Vector chuẩn 128 chiều đã được AI tính trung bình
		PasswordHash: hashedPassword,   // Lưu chuỗi đã băm
		IsFirstLogin: true,
		Status:       "ACTIVE",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 5. Lưu xuống Database
	err = s.Repo.Save(&newEmployee)
	if err != nil {
		return nil, fmt.Errorf("lỗi lưu database: %v", err)
	}

	// 6. Gửi email thông báo tài khoản
	if req.Email != nil && *req.Email != "" {
		go func() {
			err := utils.SendRegistrationEmail(*req.Email, req.Name, employeeCode, rawPassword, s.SMTPHost, s.SMTPPort, s.SMTPEmail, s.SMTPPassword)
			if err != nil {
				fmt.Printf("Lỗi gửi email đăng ký cho %s: %v\n", *req.Email, err)
			}
		}()
	}

	// 7. Hoàn thành xuất sắc
	return &dto.RegisterEmployeeResponse{
		Status:  "success",
		Message: fmt.Sprintf("Đăng ký thành công! Mã nhân viên: %s. Thông tin đã được gửi vào email.", employeeCode),
		Stats:   aiResp.Stats,
	}, nil
}

func (s *RegisterService) generateEmployeeCode() (string, error) {
	lastCode, err := s.Repo.GetLastEmployeeCode()
	if err != nil {
		return "", err
	}

	if lastCode == "" {
		return "CT000001", nil
	}

	// Giả sử format là CTxxxxxx
	numStr := lastCode[2:]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		// Nếu không parse được (ví dụ format cũ), bắt đầu lại từ đầu hoặc xử lý khác
		return "CT000001", nil
	}

	nextNum := num + 1
	return fmt.Sprintf("CT%06d", nextNum), nil
}

func (s *RegisterService) generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[num.Int64()]
	}
	return string(b)
}
