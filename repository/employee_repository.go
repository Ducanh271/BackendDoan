package repository

import (
	"database/sql"
	"duckanh/backend-doan/models"
	"encoding/json"
	"fmt"
	"time"
)

type MySQLEmployeeRepository struct {
	DB *sql.DB
}

func NewMySQLEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &MySQLEmployeeRepository{DB: db}
}

func (r *MySQLEmployeeRepository) Save(emp *models.Employee) error {
	embeddingJSON, err := json.Marshal(emp.Embedding)
	if err != nil {
		return fmt.Errorf("lỗi endcode embedding: %v", err)
	}

	// Cập nhật câu query: thêm password_hash và is_first_login (tổng cộng 13 dấu ?)
	query := `
		INSERT INTO employees (
			employee_code, name, email, phone, department_id, position_id, 
			hire_date, embedding, status, created_at, updated_at, 
			password_hash, is_first_login
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	now := time.Now()

	status := "ACTIVE"
	if emp.Status != "" {
		status = emp.Status
	}

	_, err = r.DB.Exec(query,
		emp.EmployeeCode,
		emp.Name,
		emp.Email,
		emp.Phone,
		emp.DepartmentID,
		emp.PositionID,
		emp.HireDate,
		embeddingJSON,
		status,
		now,
		now,
		emp.PasswordHash,
		emp.IsFirstLogin,
	)

	return err
}
func (r *MySQLEmployeeRepository) GetAll() ([]models.Employee, error) {
	// Chỉ lấy những nhân viên đang ACTIVE (đang làm việc) để so khớp khuôn mặt
	query := `
		SELECT id, employee_code, name, email, phone, department_id, position_id, hire_date, embedding, status, created_at, updated_at 
		FROM employees 
		WHERE status = 'ACTIVE'
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		var embeddingJSON []byte
		err := rows.Scan(
			&emp.ID,
			&emp.EmployeeCode,
			&emp.Name,
			&emp.Email,
			&emp.Phone,
			&emp.DepartmentID,
			&emp.PositionID,
			&emp.HireDate,
			&embeddingJSON,
			&emp.Status,
			&emp.CreatedAt,
			&emp.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("lỗi scan dữ liệu employee: %v", err)
		}

		if len(embeddingJSON) > 0 {
			err = json.Unmarshal(embeddingJSON, &emp.Embedding)
			if err != nil {
				return nil, fmt.Errorf("lỗi parse embedding nhân viên %s: %v", emp.Name, err)
			}
		}

		employees = append(employees, emp)
	}

	return employees, nil
}

func (r *MySQLEmployeeRepository) FindByID(id int) (*models.Employee, error) {
	query := `
		SELECT id, employee_code, name, email, phone, department_id, position_id, hire_date, embedding, status, created_at, updated_at 
		FROM employees 
		WHERE id = ? AND status = 'ACTIVE'
	`
	row := r.DB.QueryRow(query, id)

	var emp models.Employee
	var embeddingJSON []byte

	err := row.Scan(
		&emp.ID,
		&emp.EmployeeCode,
		&emp.Name,
		&emp.Email,
		&emp.Phone,
		&emp.DepartmentID,
		&emp.PositionID,
		&emp.HireDate,
		&embeddingJSON,
		&emp.Status,
		&emp.CreatedAt,
		&emp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("không tìm thấy nhân viên hoặc tài khoản bị khóa")
		}
		return nil, fmt.Errorf("lỗi scan dữ liệu employee: %v", err)
	}

	if len(embeddingJSON) > 0 {
		err = json.Unmarshal(embeddingJSON, &emp.Embedding)
		if err != nil {
			return nil, fmt.Errorf("lỗi parse embedding: %v", err)
		}
	}

	return &emp, nil
}

func (r *MySQLEmployeeRepository) LogAccess(log models.AccessLog) error {
	query := `INSERT INTO access_logs (employee_id, access_point, status, confidence_distance, failure_reason) VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, log.EmployeeID, log.AccessPoint, log.Status, log.ConfidenceDistance, log.FailureReason)
	if err != nil {
		return fmt.Errorf("lỗi lưu log truy cập: %v", err)
	}
	return nil

}
func (r *MySQLEmployeeRepository) FindByEmployeeCode(code string) (*models.Employee, error) {
	// Bổ sung otp_code và otp_expires_at vào câu SELECT
	query := `SELECT id, employee_code, name, email, password_hash, is_first_login, status, otp_code, otp_expires_at FROM employees WHERE employee_code = ?`

	var emp models.Employee
	err := r.DB.QueryRow(query, code).Scan(
		&emp.ID,
		&emp.EmployeeCode,
		&emp.Name,
		&emp.Email,
		&emp.PasswordHash,
		&emp.IsFirstLogin,
		&emp.Status,
		&emp.OTPCode,
		&emp.OTPExpiresAt,
	)

	if err != nil {
		return nil, err
	}
	return &emp, nil
}
func (r *MySQLEmployeeRepository) UpdatePasswordAndStatus(employeeID int, newHash string, isFirstLogin bool) error {
	query := `UPDATE employees SET password_hash = ?, is_first_login = ?, otp_code = NULL, otp_expires_at = NULL WHERE id = ?`
	_, err := r.DB.Exec(query, newHash, isFirstLogin, employeeID)
	return err
}

func (r *MySQLEmployeeRepository) UpdateOTP(employeeID int, otpCode string, expiresAt time.Time) error {
	query := `UPDATE employees SET otp_code = ?, otp_expires_at = ? WHERE id = ?`
	_, err := r.DB.Exec(query, otpCode, expiresAt, employeeID)
	return err
}
