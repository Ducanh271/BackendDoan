package handlers

import (
	"net/http"

	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/service"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	RegisterService *service.RegisterService
	EmployeeService *service.EmployeeService // Bổ sung thêm Service này
}

func NewEmployeeHandler(rs *service.RegisterService, es *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		RegisterService: rs,
		EmployeeService: es,
	}
}

// API: POST /api/employees/register
func (h *EmployeeHandler) Register(c *gin.Context) {
	var req dto.RegisterEmployeeRequest

	// 1. Map body JSON vào DTO (Gin sẽ tự báo lỗi nếu thiếu trường require)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	// 2. Gọi Service xử lý toàn bộ logic (AI + DB)
	res, err := h.RegisterService.RegisterEmployee(req)

	// Lỗi hệ thống (DB chết, AI sập mạng...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Lỗi do AI đánh rớt (Ảnh fake, sai khuôn mặt...) -> HTTP 422 Unprocessable Entity
	if res.Status == "failed" {
		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	// 3. Đăng ký thành công mỹ mãn -> HTTP 200 OK
	c.JSON(http.StatusOK, res)
}

func (h *EmployeeHandler) GetMyProfile(c *gin.Context) {
	empIDVal, exists := c.Get("employee_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "message": "Không xác định được danh tính"})
		return
	}

	empID := empIDVal.(int)

	profile, err := h.EmployeeService.GetProfileByID(empID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Lấy thông tin cá nhân thành công",
		"data":    profile,
	})
}
