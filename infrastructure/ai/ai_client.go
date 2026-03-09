// infrastructure/ai/ai_client.go
package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPClient struct {
	BaseURL string
}

func NewAIClient(baseURL string) AIClient {
	return &HTTPClient{BaseURL: baseURL}
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

type aiRequest struct {
	Images []string `json:"images"`
}

func (c *HTTPClient) RegisterFace(images [][]byte) (*AIRegisterResponse, error) {
	var base64Images []string
	for _, img := range images {
		base64Images = append(base64Images, encodeBase64(img))
	}

	reqBody := aiRequest{Images: base64Images}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("lỗi marshal JSON: %v", err)
	}

	apiURL := c.BaseURL + "/api/register-face"
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới AI Service: %v", err)
	}
	defer resp.Body.Close()

	var result AIRegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("lỗi parse kết quả từ AI: %v", err)
	}

	return &result, nil
}
