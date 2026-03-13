package dto

type LoginRequest struct {
	EmployeeCode string `json:"employee_code" binding:"required"`
	Password     string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IsFirstLogin bool   `json:"is_first_login"`
	Message      string `json:"message"`
}

type ChangeFirstPasswordRequest struct {
	EmployeeCode string `json:"employee_code" binding:"required"`
	OldPassword  string `json:"old_password" binding:"required"`
	NewPassword  string `json:"new_password" binding:"required,min=6"` // Bắt buộc pass mới từ 6 ký tự
}

type RequestOTPRequest struct {
	EmployeeCode string `json:"employee_code" binding:"required"`
}

type VerifyOTPAndChangePasswordRequest struct {
	EmployeeCode string `json:"employee_code" binding:"required"`
	OTPCode      string `json:"otp_code" binding:"required"`
	NewPassword  string `json:"new_password" binding:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
