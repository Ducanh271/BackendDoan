package utils

import (
	"fmt"
	"net/smtp"
)

func SendOTPEmail(toEmail, employeeName, otpCode, smtpHost, smtpPort, senderEmail, senderPass string) error {

	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)

	subject := "Subject: Mã xác thực đổi mật khẩu hệ thống DeepGuard\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf(`
		<h3>Xin chào %s,</h3>
		<p>Bạn đã yêu cầu đổi mật khẩu lần đầu cho tài khoản trên hệ thống DeepGuard.</p>
		<p>Mã xác thực (OTP) của bạn là: <b style="font-size: 20px; color: blue;">%s</b></p>
		<p>Mã này sẽ hết hạn trong 5 phút. Tuyệt đối không chia sẻ mã này cho bất kỳ ai.</p>
		<br/>
		<p>Trân trọng,<br/>Đội ngũ DeepGuard Security.</p>
	`, employeeName, otpCode)

	msg := []byte(subject + mime + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("không thể gửi email: %v", err)
	}
	return nil
}
