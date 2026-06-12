package repository

import "duckanh/backend-doan/models"

type AttendanceRepository interface {
	GetByEmployeeID(employeeID int) ([]models.AttendanceRecord, error)
	SaveAttendance(record *models.AttendanceRecord) error
	GetTodayRecord(employeeID int, date string) (*models.AttendanceRecord, error)
}
