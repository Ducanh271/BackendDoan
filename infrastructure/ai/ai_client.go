package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPClient struct {
	BaseURL   string
	client    *http.Client
	semaphore chan struct{}
}

func NewAIClient(baseURL string, maxConcurrent int) AIClient {
	return &HTTPClient{
		BaseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		semaphore: make(chan struct{}, maxConcurrent),
	}
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

type aiRequest struct {
	Images []string `json:"images"`
}

func (c *HTTPClient) RegisterFace(images [][]byte) (*AIRegisterResponse, error) {

	c.semaphore <- struct{}{}

	defer func() { <-c.semaphore }()

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

	resp, err := c.client.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
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

type aiVerifyRequest struct {
	Images []string `json:"images"`
}

func (c *HTTPClient) VerifyFace(images [][]byte) (*AIVerifyResponse, error) {
	c.semaphore <- struct{}{}
	defer func() { <-c.semaphore }()
	var base64Images []string
	for _, img := range images {
		base64Images = append(base64Images, encodeBase64(img))
	}

	reqBody := aiVerifyRequest{Images: base64Images}
	jsonData, err := json.Marshal(reqBody)

	// 3. Gọi API Python (Endpoint: /api/verify-face)
	apiURL := c.BaseURL + "/api/verify-face"
	resp, err := c.client.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("không thể kết nối tới AI Service: %v", err)
	}
	defer resp.Body.Close()

	// 4. Đọc kết quả
	var result AIVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("lỗi parse kết quả từ AI: %v", err)
	}

	return &result, nil
}
