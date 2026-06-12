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
	configHandler *handlers.ConfigHandler,
	jwtSecret string,
) {
	api := r.Group("/api")

	// ==========================================
	// 1. NHÓM AUTH (Xác thực) - CHỦ YẾU LÀ PUBLIC
	// ==========================================
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/request-otp", authHandler.RequestFirstLoginOTP)
		auth.POST("/verify-otp-change-password", authHandler.VerifyOTPAndChangePassword)
		auth.POST("/refresh-token", authHandler.RefreshToken)

		// Đăng xuất cần phải có Token mới cho đăng xuất
		auth.DELETE("/logout", middleware.RequireAuth(jwtSecret), authHandler.LogOut)
	}

	// ==========================================
	// 2. NHÓM ADMIN (Quản trị viên) - BẢO MẬT 2 LỚP
	// ==========================================
	admin := api.Group("/admin")
	// Bọc khiên: Bắt buộc đăng nhập + Bắt buộc chức vụ là ADMIN
	admin.Use(middleware.RequireAuth(jwtSecret), middleware.RequireRole("ADMIN"))
	{
		// Cấp tài khoản mới (POST /api/admin/register)
		admin.POST("/register", empHandler.Register)

		// Quản lý nhân viên
		admin.GET("/employees", empHandler.GetAllEmployees)
		admin.GET("/employees/:id", empHandler.GetEmployeeDetail)
		admin.GET("/employees/:id/attendance", attHandler.GetEmployeeAttendanceHistory)
	}

	// ==========================================
	// 3. NHÓM CẤU HÌNH (Chỉ ADMIN mới được xem/sửa)
	// ==========================================
	configs := api.Group("/configs")
	// Bọc khiên 2 lớp giống hệt nhóm Admin
	configs.Use(middleware.RequireAuth(jwtSecret), middleware.RequireRole("ADMIN"))
	{
		configs.GET("", configHandler.GetConfig)
		configs.PUT("", configHandler.UpdateConfig)
		configs.POST("/wifis", configHandler.AddWifi)
		configs.DELETE("/wifis/:id", configHandler.DeleteWifi)
	}

	// ==========================================
	// 4. NHÓM NHÂN VIÊN (Dành cho mọi User đã đăng nhập)
	// ==========================================
	employees := api.Group("/employees")
	employees.Use(middleware.RequireAuth(jwtSecret))
	{
		employees.GET("/me", empHandler.GetMyProfile)
	}

	// ==========================================
	// 5. NHÓM ATTENDANCE (Điểm danh)
	// ==========================================
	attendance := api.Group("/attendance")
	{
		// [Public] Điểm danh qua Tablet ở cửa
		attendance.POST("/identify", attHandler.Identify)

		// [Protected] Điểm danh qua Mobile App cá nhân
		attendance.POST("/mobile-check-in", middleware.RequireAuth(jwtSecret), attHandler.MobileCheckIn)

		// [Protected] Xem lịch sử điểm danh cá nhân
		attendance.GET("/history", middleware.RequireAuth(jwtSecret), attHandler.GetHistory)
	}
}
