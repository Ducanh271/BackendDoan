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

func (r *MySQLEmployeeRepository) Save(emp models.Employee) error {
	embeddingJSON, err := json.Marshal(emp.Embedding)
	if err != nil {
		return fmt.Errorf("lỗi endcode embedding: %v", err)
	}
	query := `
		INSERT INTO employees (employee_code, name, email, phone, department_id, position_id, hire_date, embedding, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
