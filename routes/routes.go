package routes

import (
	"duckanh/backend-doan/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes thiết lập toàn bộ các API endpoint cho hệ thống
func SetupRoutes(
	r *gin.Engine,
	empHandler *handlers.EmployeeHandler,
	attHandler *handlers.AttendanceHandler,
	authHandler *handlers.AuthHandler,
) {
	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		// api.POST("/change-first-password", authHandler.ChangeFirstPassword)
		api.POST("/request-otp", authHandler.RequestFirstLoginOTP)
		api.POST("/verify-otp-change-password", authHandler.VerifyOTPAndChangePassword)

		employees := api.Group("/employees")
		{
			employees.POST("/register", empHandler.Register)
			// Tương lai: employees.GET("/", empHandler.GetAll)
			// Tương lai: employees.GET("/:id", empHandler.GetByID)
		}

		attendance := api.Group("/attendance")
		{
			attendance.POST("/identify", attHandler.Identify)
			// Sau này làm Android: attendance.POST("/check-in", AuthMiddleware, attHandler.CheckIn)
		}
	}
}
