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
) {
	api := r.Group("/api")
	{
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
