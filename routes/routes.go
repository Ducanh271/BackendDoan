package routes

import (
	"duckanh/backend-doan/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes thiết lập toàn bộ các API endpoint cho hệ thống
func SetupRoutes(r *gin.Engine, empHandler *handlers.EmployeeHandler) {
	// Nhóm các API lại với prefix là /api
	api := r.Group("/api")
	{
		// Nhóm API liên quan đến Employee
		employees := api.Group("/employees")
		{
			employees.POST("/register", empHandler.Register)
			// Tương lai: employees.GET("/", empHandler.GetAll)
			// Tương lai: employees.GET("/:id", empHandler.GetByID)
		}

		// Nhóm API liên quan đến Attendance (chuẩn bị sẵn chỗ)
		// attendance := api.Group("/attendance")
		// {
		// 	attendance.POST("/verify", ...)
		// }
	}
}
