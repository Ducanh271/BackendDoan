package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Không thể kết nối database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("DB không phản hồi (Ping failed): %v", err)
	}
	fmt.Println("Đã kết nối thành công tới database (0.0)")
	return db
}
