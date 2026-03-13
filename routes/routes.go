package routes

import (
	"duckanh/backend-doan/handlers"
	"duckanh/backend-doan/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes thiết lập toàn bộ các API endpoint cho hệ thống
func SetupRoutes(
	r *gin.Engine,
	empHandler *handlers.EmployeeHandler,
	attHandler *handlers.AttendanceHandler,
	authHandler *handlers.AuthHandler,
	jwtSecret string,
) {
	api := r.Group("/api")

	// ==========================================
	// 1. NHÓM AUTH (Xác thực & Cấp quyền) - KHÔNG CẦN TOKEN
	// ==========================================
	auth := api.Group("/auth")
	{
		auth.POST("/register", empHandler.Register) // Chuyển Register về đúng "nhà" của nó
		auth.POST("/login", authHandler.Login)
		auth.POST("/request-otp", authHandler.RequestFirstLoginOTP)
		auth.POST("/verify-otp-change-password", authHandler.VerifyOTPAndChangePassword)
		auth.POST("/refresh-token", authHandler.RefreshToken)
	}

	// ==========================================
	// 2. NHÓM EMPLOYEES (Quản lý nhân viên) - BẮT BUỘC CÓ TOKEN
	// ==========================================
	employees := api.Group("/employees")
	employees.Use(middleware.RequireAuth(jwtSecret)) // Bọc khiên bảo vệ cho TOÀN BỘ group này
	{
		// Đường dẫn thực tế: GET /api/employees/me
		employees.GET("/me", empHandler.GetMyProfile)

		// Sau này cậu viết thêm các API như:
		// employees.PUT("/:id", empHandler.UpdateEmployee)
		// employees.DELETE("/:id", empHandler.DeleteEmployee)
		// ... tất cả sẽ tự động được bảo vệ!
	}

	// ==========================================
	// 3. NHÓM ATTENDANCE (Điểm danh) - HỖN HỢP
	// ==========================================
	attendance := api.Group("/attendance")
	{
		// [Public] API Điểm danh qua Camera chung để ở sảnh (Không cần Token)
		attendance.POST("/identify", attHandler.Identify)

		// 💡 [Protected] API Điểm danh cá nhân trên Mobile App (Cần Token)
		// Kẹp middleware trực tiếp vào từng route cụ thể trong group hỗn hợp này:
		// attendance.POST("/mobile-check-in", middleware.RequireAuth(jwtSecret), attHandler.MobileCheckIn)
	}
}
