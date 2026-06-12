package repository

import (
	"database/sql"
	"duckanh/backend-doan/models"
	"fmt"
)

type MySQLAttendanceRepository struct {
	DB *sql.DB
}

func NewMySQLAttendanceRepository(db *sql.DB) AttendanceRepository {
	return &MySQLAttendanceRepository{DB: db}
}

func (r *MySQLAttendanceRepository) GetByEmployeeID(employeeID int) ([]models.AttendanceRecord, error) {
	query := `
		SELECT id, employee_id, date, check_in, check_out, location, method, created_at
		FROM attendance_records
		WHERE employee_id = ?
		ORDER BY date DESC, check_in DESC
	`
	rows, err := r.DB.Query(query, employeeID)
	if err != nil {
		return nil, fmt.Errorf("lỗi truy vấn lịch sử điểm danh: %v", err)
	}
	defer rows.Close()

	var records []models.AttendanceRecord
	for rows.Next() {
		var rec models.AttendanceRecord
		err := rows.Scan(
			&rec.ID,
			&rec.EmployeeID,
			&rec.Date,
			&rec.CheckIn,
			&rec.CheckOut,
			&rec.Location,
			&rec.Method,
			&rec.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("lỗi scan dữ liệu điểm danh: %v", err)
		}
		records = append(records, rec)
	}

	return records, nil
}

func (r *MySQLAttendanceRepository) GetTodayRecord(employeeID int, date string) (*models.AttendanceRecord, error) {
	query := `
		SELECT id, employee_id, date, check_in, check_out, location, method, created_at
		FROM attendance_records
		WHERE employee_id = ? AND date = ?
		LIMIT 1
	`
	var rec models.AttendanceRecord
	err := r.DB.QueryRow(query, employeeID, date).Scan(
		&rec.ID,
		&rec.EmployeeID,
		&rec.Date,
		&rec.CheckIn,
		&rec.CheckOut,
		&rec.Location,
		&rec.Method,
		&rec.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r *MySQLAttendanceRepository) SaveAttendance(rec *models.AttendanceRecord) error {
	if rec.ID == 0 {
		// Tạo mới (Check-in)
		query := `
			INSERT INTO attendance_records (employee_id, date, check_in, location, method)
			VALUES (?, ?, ?, ?, ?)
		`
		_, err := r.DB.Exec(query, rec.EmployeeID, rec.Date, rec.CheckIn, rec.Location, rec.Method)
		return err
	} else {
		// Cập nhật (Check-out)
		query := `
			UPDATE attendance_records
			SET check_out = ?
			WHERE id = ?
		`
		_, err := r.DB.Exec(query, rec.CheckOut, rec.ID)
		return err
	}
}
