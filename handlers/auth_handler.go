// handlers/auth_handler.go
package handlers

import (
	"net/http"

	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(as *service.AuthService) *AuthHandler {
	return &AuthHandler{AuthService: as}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "Dữ liệu không hợp lệ: vui lòng nhập mã nhân viên và mật khẩu",
		})
		return
	}

	res, err := h.AuthService.Login(req)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   res,
	})
}

// func (h *AuthHandler) ChangeFirstPassword(c *gin.Context) {
// 	var req dto.ChangeFirstPasswordRequest
//
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  "failed",
// 			"message": "Dữ liệu không hợp lệ, vui lòng kiểm tra lại (mật khẩu mới phải từ 6 ký tự)",
// 		})
// 		return
// 	}
//
// 	err := h.AuthService.ChangeFirstPassword(req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  "failed",
// 			"message": err.Error(),
// 		})
// 		return
// 	}
//
// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "success",
// 		"message": "Đổi mật khẩu thành công! Vui lòng đăng nhập lại với mật khẩu mới.",
// 	})
// }

func (h *AuthHandler) RequestFirstLoginOTP(c *gin.Context) {
	var req dto.RequestOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Vui lòng cung cấp mã nhân viên"})
		return
	}

	if err := h.AuthService.RequestFirstLoginOTP(req.EmployeeCode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mã xác thực đã được gửi đến email của bạn. Vui lòng kiểm tra hộp thư.",
	})
}

func (h *AuthHandler) VerifyOTPAndChangePassword(c *gin.Context) {
	var req dto.VerifyOTPAndChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": "Dữ liệu không hợp lệ (mật khẩu mới phải từ 6 ký tự)"})
		return
	}

	if err := h.AuthService.VerifyOTPAndChangePassword(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Đổi mật khẩu thành công! Tài khoản đã được kích hoạt.",
	})
}
