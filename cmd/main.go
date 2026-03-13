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
	"duckanh/backend-doan/routes"
	"duckanh/backend-doan/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("🚀 Khởi động hệ thống Backend...")

	// 1. Load cấu hình
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println("⚠️ Cảnh báo:", err)
	}

	// 2. Kết nối DB
	db := database.Connect(cfg.DB.URL)
	defer db.Close()

	// 3. Cấu hình AI
	aiURL := os.Getenv("AI_SERVICE_URL")
	if aiURL == "" {
		aiURL = "http://127.0.0.1:8000" // Cổng của Python FastAPI
	}

	// 4. Khởi tạo các "Khối Lego" (Dependency Injection)
	aiClient := ai.NewAIClient(aiURL)

	empRepo := repository.NewMySQLEmployeeRepository(db)
	tokenRepo := repository.NewMySQLRefreshTokenRepository(db)

	regService := service.NewRegisterService(empRepo, aiClient)
	attService := service.NewAttendanceService(empRepo, aiClient)
	authService := service.NewAuthService(empRepo, tokenRepo, cfg.JWTSecret, cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPEmail, cfg.SMTPPassword)
	empService := service.NewEmployeeService(empRepo)

	empHandler := handlers.NewEmployeeHandler(regService, empService)
	attHandler := handlers.NewAttendanceHandler(attService)
	authHandler := handlers.NewAuthHandler(authService)

	// 5. Thiết lập Router (Gin)
	r := gin.Default()

	// === BẮT BUỘC PHẢI CÓ ĐOẠN NÀY ĐỂ REACT GỌI ĐƯỢC API ===
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true // Cho phép mọi Frontend gọi vào
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))
	// =======================================================

	routes.SetupRoutes(r, empHandler, attHandler, authHandler, cfg.JWTSecret)
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
