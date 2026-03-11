package handlers

import (
	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AttendanceHandler struct {
	AttendanceService *service.AttendanceService
}

func NewAttendanceHandler(as *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{AttendanceService: as}
}

// API: POST /api/attendance/identify (Dành cho Web 1:N)
func (h *AttendanceHandler) Identify(c *gin.Context) {
	var req dto.VerifyAttendanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Dữ liệu không hợp lệ: " + err.Error(),
		})
		return
	}

	// Gọi luồng 1:N
	emp, dist, aiResp, err := h.AttendanceService.IdentifyFace(req.Images)

	// TRƯỜNG HỢP 1: NẾU BỊ LỖI (AI chê fake, không thấy mặt, hoặc không khớp ai trong DB)
	if err != nil {
		// Gọi qua Service để lưu log âm thầm
		_ = h.AttendanceService.LogAccessEvent(nil, "Cổng chính", "DENIED", dist, err.Error())

		res := dto.VerifyAttendanceResponse{
			Status:  "failed",
			Message: err.Error(),
		}
		// Đính kèm Stats của AI (nếu có) để Frontend hiện lý do rớt
		if aiResp != nil {
			res.Stats = aiResp.Stats
		}
		c.JSON(http.StatusUnprocessableEntity, res)
		return
	}

	// TRƯỜNG HỢP 2: NẾU NHẬN DIỆN THÀNH CÔNG

	// Gọi qua Service để lưu log âm thầm
	_ = h.AttendanceService.LogAccessEvent(&emp.ID, "Cổng chính", "GRANTED", dist, "")

	// Trả kết quả thành công về Frontend
	c.JSON(http.StatusOK, dto.VerifyAttendanceResponse{
		Status:       "success",
		Message:      "Nhận diện thành công!",
		EmployeeName: emp.Name,
		EmployeeCode: emp.EmployeeCode,
		Distance:     dist,
		Stats:        aiResp.Stats,
	})
}
