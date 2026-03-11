package service

import (
	"duckanh/backend-doan/infrastructure/ai"
	"duckanh/backend-doan/models"
	"duckanh/backend-doan/repository"
	"duckanh/backend-doan/utils"
	"errors"
	"fmt"
	"math"
)

const FaceMatchThreshold = 0.45

type AttendanceService struct {
	EmpRepo repository.EmployeeRepository
	AI      ai.AIClient
}

func NewAttendanceService(repo repository.EmployeeRepository, aiClient ai.AIClient) *AttendanceService {
	return &AttendanceService{
		EmpRepo: repo,
		AI:      aiClient,
	}
}

// ==========================================
// 1:1 - CHECK-IN TỪ APP ANDROID (Có ID từ Token)
// ==========================================
func (s *AttendanceService) CheckIn(employeeID int, imagesBase64 []string) (*models.Employee, float64, *ai.AIVerifyResponse, error) {
	emp, err := s.EmpRepo.FindByID(employeeID)
	if err != nil {
		return nil, 0, nil, err
	}

	imgBytes, err := utils.DecodeBase64Images(imagesBase64)
	if err != nil {
		return nil, 0, nil, err
	}

	aiResp, err := s.AI.VerifyFace(imgBytes)
	if err != nil {
		return nil, 0, nil, err
	}
	if aiResp.Status != "success" {
		return nil, 0, aiResp, fmt.Errorf("AI từ chối: %s", aiResp.Error)
	}

	// SỬ DỤNG UTILS.EUCLIDEANDISTANCE
	dist := utils.EuclideanDistance(emp.Embedding, aiResp.Embedding)
	if dist == math.MaxFloat64 {
		return nil, 0, aiResp, errors.New("kích thước vector (embedding) không hợp lệ")
	}

	if dist > FaceMatchThreshold {
		return nil, dist, aiResp, fmt.Errorf("khuôn mặt không khớp (sai số: %.2f)", dist)
	}

	return emp, dist, aiResp, nil
}

// ==========================================
// 1:N - IDENTIFY TỪ WEB QUẢN TRỊ (Không có ID)
// ==========================================
func (s *AttendanceService) IdentifyFace(imagesBase64 []string) (*models.Employee, float64, *ai.AIVerifyResponse, error) {
	imgBytes, err := utils.DecodeBase64Images(imagesBase64)
	if err != nil {
		return nil, 0, nil, err
	}

	aiResp, err := s.AI.VerifyFace(imgBytes)
	if err != nil {
		return nil, 0, nil, err
	}
	if aiResp.Status != "success" {
		return nil, 0, aiResp, fmt.Errorf("AI từ chối: %s", aiResp.Error)
	}

	employees, err := s.EmpRepo.GetAll()
	if err != nil {
		return nil, 0, aiResp, fmt.Errorf("lỗi lấy danh sách NV: %v", err)
	}

	var bestMatch *models.Employee
	minDist := math.MaxFloat64

	for i := range employees {
		// SỬ DỤNG UTILS.EUCLIDEANDISTANCE
		dist := utils.EuclideanDistance(employees[i].Embedding, aiResp.Embedding)

		// Bỏ qua nếu lỗi chiều dài vector, ngược lại kiểm tra xem có phải khoảng cách nhỏ nhất không
		if dist != math.MaxFloat64 && dist < minDist {
			minDist = dist
			bestMatch = &employees[i]
		}
	}

	if bestMatch == nil || minDist > FaceMatchThreshold {
		return nil, minDist, aiResp, errors.New("không nhận diện được khuôn mặt này trong hệ thống")
	}

	return bestMatch, minDist, aiResp, nil
}

func (s *AttendanceService) LogAccessEvent(employeeID *int, accessPoint string, status string, distance float64, failureReason string) error {
	accessLog := models.AccessLog{
		EmployeeID:         employeeID,
		AccessPoint:        accessPoint,
		Status:             status,
		ConfidenceDistance: distance,
		FailureReason:      failureReason,
	}

	return s.EmpRepo.LogAccess(accessLog)
}
