package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type DBConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string
	URL  string
}
type Config struct {
	DB DBConfig
}

func LoadConfig() (Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Không tìm thấy file .env, dùng biến môi trường có sẵn")
	}
	var cfg Config
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Pass = os.Getenv("DB_PASS")
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.URL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)
	return cfg, nil
}
