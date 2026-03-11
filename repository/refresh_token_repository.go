package repository

import (
	"database/sql"
	"duckanh/backend-doan/models"
	"fmt"
)

type MySQLRefreshTokenRepository struct {
	DB *sql.DB
}

func NewMySQLRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return &MySQLRefreshTokenRepository{DB: db}
}

func (r *MySQLRefreshTokenRepository) Save(token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (employee_id, token, expires_at)
		VALUES (?, ?, ?)
	`
	result, err := r.DB.Exec(query, token.EmployeeID, token.Token, token.ExpiresAt)
	if err != nil {
		return fmt.Errorf("lỗi lưu refresh token: %v", err)
	}

	id, _ := result.LastInsertId()
	token.ID = int(id)
	return nil
}

func (r *MySQLRefreshTokenRepository) FindByToken(tokenString string) (*models.RefreshToken, error) {
	query := `
		SELECT id, employee_id, token, expires_at, created_at 
		FROM refresh_tokens 
		WHERE token = ?
	`
	var t models.RefreshToken
	err := r.DB.QueryRow(query, tokenString).Scan(
		&t.ID,
		&t.EmployeeID,
		&t.Token,
		&t.ExpiresAt,
		&t.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("không tìm thấy token hoặc token đã bị thu hồi")
		}
		return nil, fmt.Errorf("lỗi truy vấn token: %v", err)
	}

	return &t, nil
}

func (r *MySQLRefreshTokenRepository) DeleteByToken(tokenString string) error {
	query := `DELETE FROM refresh_tokens WHERE token = ?`
	_, err := r.DB.Exec(query, tokenString)
	return err
}

func (r *MySQLRefreshTokenRepository) DeleteAllByEmployeeID(employeeID int) error {
	query := `DELETE FROM refresh_tokens WHERE employee_id = ?`
	_, err := r.DB.Exec(query, employeeID)
	return err
}
