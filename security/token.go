package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(employeeID int, secretKey string, role string) (string, error) {
	claims := jwt.MapClaims{
		"employee_id": employeeID,
		"role":        role,
		"exp":         time.Now().Add(30 * time.Minute).Unix(),
		"iat":         time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

// 2. SINH REFRESH TOKEN
func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("lỗi sinh chuỗi ngẫu nhiên: %v", err)
	}

	return hex.EncodeToString(bytes), nil
}
