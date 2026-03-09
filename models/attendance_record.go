package models

import "time"

type AttendanceRecord struct {
	ID         int        `json:"id"`
	EmployeeID int        `json:"employee_id"`
	Date       string     `json:"date"` // Lưu dạng "YYYY-MM-DD"
	CheckIn    time.Time  `json:"check_in"`
	CheckOut   *time.Time `json:"check_out"` // Dùng con trỏ vì ban đầu check_out sẽ là NULL
	Location   string     `json:"location"`
	Method     string     `json:"method"`
	CreatedAt  time.Time  `json:"created_at"`
}
