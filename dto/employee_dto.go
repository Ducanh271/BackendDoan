package dto

type RegisterEmployeeRequest struct {
	EmployeeCode string   `json:"employee_code" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Email        *string  `json:"email"`
	Phone        *string  `json:"phone"`
	DepartmentID *int     `json:"department_id"`
	PositionID   *int     `json:"position_id"`
	Images       []string `json:"images" binding:"required"` // Danh sách 5 ảnh Base64
}

type RegisterEmployeeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	// Có thể thêm các thống kê từ AI vào đây để hiển thị lên App nếu muốn
	Stats interface{} `json:"stats,omitempty"`
}
