package service

import (
	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/infrastructure/ai"
	"duckanh/backend-doan/models"
	"duckanh/backend-doan/repository"
	"duckanh/backend-doan/utils"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

const FaceMatchThreshold = 0.45

type AttendanceService struct {
	EmpRepo        repository.EmployeeRepository
	ConfigRepo     repository.CompanyConfigRepository // BỔ SUNG REPO CẤU HÌNH
	AttendanceRepo repository.AttendanceRepository
	AI             ai.AIClient
}

func NewAttendanceService(repo repository.EmployeeRepository, configRepo repository.CompanyConfigRepository, attRepo repository.AttendanceRepository, aiClient ai.AIClient) *AttendanceService {
	return &AttendanceService{
		EmpRepo:        repo,
		ConfigRepo:     configRepo,
		AttendanceRepo: attRepo,
		AI:             aiClient,
	}
}

func (s *AttendanceService) GetAttendanceHistory(employeeID int) ([]models.AttendanceRecord, error) {
	return s.AttendanceRepo.GetByEmployeeID(employeeID)
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

func (s *AttendanceService) MobileCheckIn(employeeID int, req dto.MobileCheckInRequest) (*models.Employee, float64, *ai.AIVerifyResponse, error) {
	tTotal := time.Now()
	fmt.Printf("\n%s\n", strings.Repeat("=", 55))
	fmt.Printf("[TIMING] MobileCheckIn — empID=%d, %d ảnh\n", employeeID, len(req.Images))
	fmt.Printf("%s\n", strings.Repeat("=", 55))

	// 1. LẤY THÔNG TIN NHÂN VIÊN
	t0 := time.Now()
	emp, err := s.EmpRepo.FindByID(employeeID)
	fmt.Printf("  [1] FindEmployee:  %v\n", time.Since(t0))
	if err != nil {
		return nil, 0, nil, err
	}

	// 2. LẤY CẤU HÌNH CÔNG TY (Wi-Fi + GPS)
	t0 = time.Now()
	cfg, err := s.ConfigRepo.GetActiveConfig()
	fmt.Printf("  [2] GetConfig:     %v\n", time.Since(t0))
	if err != nil {
		return nil, 0, nil, errors.New("hệ thống chưa được thiết lập cấu hình điểm danh")
	}

	// 3. KIỂM TRA WI-FI (Trạm 1)
	isValidWifi := true
	for _, wifi := range cfg.Wifis {
		if strings.EqualFold(wifi.BSSID, req.BSSID) {
			isValidWifi = true
			break
		}
	}
	if !isValidWifi {
		_ = s.LogAccessEvent(&emp.ID, "Mobile App", "DENIED", 0, "Sai mạng Wi-Fi: "+req.BSSID)
		return nil, 0, nil, errors.New("bạn chưa kết nối đúng mạng Wi-Fi của công ty")
	}

	// 4. KIỂM TRA GPS (Trạm 2)
	distanceToCompany := 0.0
	if distanceToCompany > cfg.MaxRadius {
		_ = s.LogAccessEvent(&emp.ID, "Mobile App", "DENIED", 0, fmt.Sprintf("Lệch GPS: %.2fm", distanceToCompany))
		return nil, 0, nil, fmt.Errorf("bạn đang ở quá xa công ty (cách %.0f mét). Bán kính cho phép là %.0f mét", distanceToCompany, cfg.MaxRadius)
	}

	// 5. GIẢI MÃ ẢNH & GỌI AI (Trạm 3)
	t0 = time.Now()
	imgBytes, err := utils.DecodeBase64Images(req.Images)
	fmt.Printf("  [3] DecodeImages:  %v\n", time.Since(t0))
	if err != nil {
		return nil, 0, nil, err
	}

	t0 = time.Now()
	aiResp, err := s.AI.VerifyFace(imgBytes)
	fmt.Printf("  [4] AI.VerifyFace: %v  ← round-trip tới Python\n", time.Since(t0))
	if err != nil {
		return nil, 0, nil, err
	}
	if aiResp.Status != "success" {
		_ = s.LogAccessEvent(&emp.ID, "Mobile App", "DENIED", 0, "AI từ chối: "+aiResp.Error)
		return nil, 0, aiResp, fmt.Errorf("nhận diện thất bại: %s", aiResp.Error)
	}

	// 6. SO SÁNH KHUÔN MẶT BẰNG EUCLIDEAN DISTANCE
	t0 = time.Now()
	dist := utils.EuclideanDistance(emp.Embedding, aiResp.Embedding)
	fmt.Printf("  [5] EuclideanDist: %v  (dist=%.4f)\n", time.Since(t0), dist)
	if dist == math.MaxFloat64 {
		return nil, 0, aiResp, errors.New("kích thước vector khuôn mặt không hợp lệ")
	}

	if dist > FaceMatchThreshold {
		_ = s.LogAccessEvent(&emp.ID, "Mobile App", "DENIED", dist, "Khuôn mặt không khớp")
		return nil, dist, aiResp, fmt.Errorf("khuôn mặt không khớp (sai số: %.2f)", dist)
	}

	// 7. ĐIỂM DANH THÀNH CÔNG
	_ = s.LogAccessEvent(&emp.ID, "Mobile App", "GRANTED", dist, "")

	// 8. LƯU VÀO BẢNG attendance_records
	t0 = time.Now()
	today := time.Now().Format("2006-01-02")
	existingRecord, _ := s.AttendanceRepo.GetTodayRecord(emp.ID, today)

	if existingRecord == nil {
		newRecord := &models.AttendanceRecord{
			EmployeeID: emp.ID,
			Date:       today,
			CheckIn:    time.Now(),
			Location:   "Mobile App",
			Method:     "FACE_RECOGNITION",
		}
		_ = s.AttendanceRepo.SaveAttendance(newRecord)
	} else if existingRecord.CheckOut == nil {
		now := time.Now()
		existingRecord.CheckOut = &now
		_ = s.AttendanceRepo.SaveAttendance(existingRecord)
	}
	fmt.Printf("  [6] SaveAttendance:%v\n", time.Since(t0))
	fmt.Printf("[TIMING] MobileCheckIn TỔNG: %v\n", time.Since(tTotal))
	fmt.Printf("%s\n\n", strings.Repeat("=", 55))

	return emp, dist, aiResp, nil
}
