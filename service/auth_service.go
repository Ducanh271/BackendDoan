package service

import (
	"duckanh/backend-doan/utils"
	"errors"
	"fmt"
	"time"

	"duckanh/backend-doan/dto"
	"duckanh/backend-doan/models"
	"duckanh/backend-doan/repository"
	"duckanh/backend-doan/security" // Dùng package security cậu vừa tách ra
)

type AuthService struct {
	EmpRepo      repository.EmployeeRepository
	TokenRepo    repository.RefreshTokenRepository
	JWTSecret    string
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string
}

func NewAuthService(
	empRepo repository.EmployeeRepository,
	tokenRepo repository.RefreshTokenRepository,
	secret string,
	smtpHost, smtpPort, smtpEmail, smtpPass string,
) *AuthService {
	return &AuthService{
		EmpRepo:      empRepo,
		TokenRepo:    tokenRepo,
		JWTSecret:    secret,
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPEmail:    smtpEmail,
		SMTPPassword: smtpPass,
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	emp, err := s.EmpRepo.FindByEmployeeCode(req.EmployeeCode)
	if err != nil {
		return nil, errors.New("sai mã nhân viên hoặc mật khẩu")
	}

	if emp.Status != "ACTIVE" {
		return nil, errors.New("tài khoản của bạn đã bị khóa")
	}

	if !security.CheckPasswordHash(req.Password, emp.PasswordHash) {
		return nil, errors.New("sai mã nhân viên hoặc mật khẩu")
	}
	if emp.IsFirstLogin {
		return &dto.LoginResponse{
			AccessToken:  "", // TUYỆT ĐỐI KHÔNG CẤP
			RefreshToken: "", // TUYỆT ĐỐI KHÔNG CẤP
			IsFirstLogin: true,
			Message:      "Yêu cầu đổi mật khẩu mặc định trước khi truy cập hệ thống",
		}, nil
	}

	accessToken, err := security.GenerateAccessToken(emp.ID, s.JWTSecret)
	if err != nil {
		return nil, errors.New("lỗi hệ thống khi tạo access token")
	}

	refreshTokenStr, err := security.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("lỗi hệ thống khi tạo refresh token")
	}

	rt := models.RefreshToken{
		EmployeeID: emp.ID,
		Token:      refreshTokenStr,
		ExpiresAt:  time.Now().Add(30 * 24 * time.Hour),
	}

	if err := s.TokenRepo.Save(&rt); err != nil {
		return nil, errors.New("lỗi hệ thống khi lưu phiên đăng nhập")
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		IsFirstLogin: false,
		Message:      "Đăng nhập thành công",
	}, nil
}

//	func (s *AuthService) ChangeFirstPassword(req dto.ChangeFirstPasswordRequest) error {
//		emp, err := s.EmpRepo.FindByEmployeeCode(req.EmployeeCode)
//		if err != nil {
//			return errors.New("tài khoản không tồn tại")
//		}
//
//		if !emp.IsFirstLogin {
//			return errors.New("tài khoản này đã qua lần đăng nhập đầu tiên, vui lòng dùng chức năng quên mật khẩu nếu cần")
//		}
//
//		if !security.CheckPasswordHash(req.OldPassword, emp.PasswordHash) {
//			return errors.New("mật khẩu cũ không chính xác")
//		}
//
//		if security.CheckPasswordHash(req.NewPassword, emp.PasswordHash) {
//			return errors.New("mật khẩu mới không được trùng với mật khẩu mặc định")
//		}
//
//		newHash, err := security.HashPassword(req.NewPassword)
//		if err != nil {
//			return errors.New("lỗi mã hóa mật khẩu hệ thống")
//		}
//
//		err = s.EmpRepo.UpdatePasswordAndStatus(emp.ID, newHash, false)
//		if err != nil {
//			return errors.New("lỗi khi cập nhật mật khẩu xuống cơ sở dữ liệu")
//		}
//
//		return nil
//	}
func (s *AuthService) RequestFirstLoginOTP(employeeCode string) error {
	emp, err := s.EmpRepo.FindByEmployeeCode(employeeCode)
	if err != nil {
		return errors.New("tài khoản không tồn tại")
	}

	if !emp.IsFirstLogin {
		return errors.New("tài khoản này đã được kích hoạt, không thể gửi lại OTP đổi mật khẩu lần đầu")
	}

	if *emp.Email == "" {
		return errors.New("tài khoản chưa được liên kết email, vui lòng liên hệ Admin")
	}

	otpCode, err := security.Generate6DigitOTP()
	if err != nil {
		return errors.New("lỗi hệ thống khi sinh mã xác thực")
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	err = s.EmpRepo.UpdateOTP(emp.ID, otpCode, expiresAt)
	if err != nil {
		return errors.New("lỗi hệ thống khi lưu mã xác thực")
	}

	go func() {
		errMail := utils.SendOTPEmail(
			*emp.Email, emp.Name, otpCode,
			s.SMTPHost, s.SMTPPort, s.SMTPEmail, s.SMTPPassword,
		)
		if errMail != nil {
			fmt.Printf(" Lỗi gửi email OTP cho %s: %v\n", emp.Email, errMail)
		} else {
			fmt.Printf(" Đã gửi email OTP thành công cho %s\n", emp.Email)
		}
	}()

	return nil
}

// service/auth_service.go

func (s *AuthService) VerifyOTPAndChangePassword(req dto.VerifyOTPAndChangePasswordRequest) error {
	emp, err := s.EmpRepo.FindByEmployeeCode(req.EmployeeCode)
	if err != nil {
		return errors.New("tài khoản không tồn tại")
	}

	if !emp.IsFirstLogin {
		return errors.New("tài khoản này đã được kích hoạt, không thể đổi mật khẩu theo luồng này")
	}

	// 3. KIỂM TRA OTP
	if emp.OTPCode == nil || *emp.OTPCode != req.OTPCode {
		return errors.New("mã xác thực (OTP) không chính xác")
	}

	if emp.OTPExpiresAt == nil || emp.OTPExpiresAt.Before(time.Now()) {
		return errors.New("mã xác thực đã hết hạn, vui lòng yêu cầu gửi lại mã mới")
	}

	if security.CheckPasswordHash(req.NewPassword, emp.PasswordHash) {
		return errors.New("mật khẩu mới không được giống với mật khẩu mặc định")
	}

	newHash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("lỗi mã hóa mật khẩu hệ thống")
	}

	err = s.EmpRepo.UpdatePasswordAndStatus(emp.ID, newHash, false)
	if err != nil {
		return errors.New("lỗi khi cập nhật mật khẩu xuống cơ sở dữ liệu")
	}

	return nil
}
