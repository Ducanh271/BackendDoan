// cmd/main.go
package main

import (
	"fmt"
	"log"
	"os"

	"duckanh/backend-doan/config"
	"duckanh/backend-doan/database"
	"duckanh/backend-doan/handlers"
	"duckanh/backend-doan/infrastructure/ai"
	"duckanh/backend-doan/repository"
	"duckanh/backend-doan/service"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 Khởi động hệ thống Backend...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("⚠️ Cảnh báo:", err)
	}

	db := database.Connect(cfg.DB.URL)
	defer db.Close()

	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		aiURL = "http://127.0.0.1:8000" // Cổng của Python FastAPI
	}

	// 4. Khởi tạo các "Khối Lego" (Dependency Injection)
	// - Tầng Infrastructure & Repository
	aiClient := ai.NewAIClient(aiURL)
	empRepo := repository.NewMySQLEmployeeRepository(db)

	// - Tầng Service
	regService := service.NewRegisterService(empRepo, aiClient)

	// - Tầng Handler
	empHandler := handlers.NewEmployeeHandler(regService)

	// 5. Thiết lập Router (Gin)
	r := gin.Default()

	// Định tuyến các API
	api := r.Group("/api")
	{
		api.POST("/employees/register", empHandler.Register)
		// Các API tiếp theo (ví dụ: điểm danh, mở file) sẽ thêm vào đây sau
	}

	// 6. Chạy Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("✅ Server Backend đang chạy tại http://localhost:%s\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Lỗi khởi chạy server:", err)
	}
}
