package utils

import (
	"encoding/base64"
	"fmt"
	"math"
	"strings"
)

func CosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) {
		return 0
	}

	var dot, normA, normB float64

	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
func EuclideanDistance(a, b []float64) float64 {
	if len(a) != len(b) {
		return math.MaxFloat64 // Trả về vô cùng nếu lỗi độ dài
	}

	var sum float64
	for i := 0; i < len(a); i++ {
		diff := a[i] - b[i]
		sum += diff * diff
	}

	return math.Sqrt(sum)
}

func DecodeBase64Images(b64Strings []string) ([][]byte, error) {
	var imageBytes [][]byte
	for _, b64Str := range b64Strings {
		cleanBase64 := b64Str
		if idx := strings.Index(cleanBase64, ","); idx != -1 {
			cleanBase64 = cleanBase64[idx+1:]
		}
		decoded, err := base64.StdEncoding.DecodeString(cleanBase64)
		if err != nil {
			return nil, fmt.Errorf("lỗi giải mã ảnh base64: %v", err)
		}
		imageBytes = append(imageBytes, decoded)
	}
	return imageBytes, nil
}
