package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

const (
	apiURL      = "http://localhost:8080/api/attendance/identify"
	imageDir    = "../test_images"
	numRequests = 20
)

type VerifyAttendanceRequest struct {
	Images []string `json:"images"`
}

func main() {
	fmt.Printf("Bắt đầu test đồng thời %d request tới %s...\n", numRequests, apiURL)

	// 1. Đọc tất cả ảnh từ folder test_images
	files, err := ioutil.ReadDir(imageDir)
	if err != nil {
		fmt.Printf("Lỗi đọc folder ảnh: %v\n", err)
		return
	}

	var base64Images []string
	for _, f := range files {
		if !f.IsDir() {
			ext := filepath.Ext(f.Name())
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				path := filepath.Join(imageDir, f.Name())
				data, err := ioutil.ReadFile(path)
				if err != nil {
					fmt.Printf("Lỗi đọc file %s: %v\n", f.Name(), err)
					continue
				}
				b64 := base64.StdEncoding.EncodeToString(data)
				base64Images = append(base64Images, b64)
			}
		}
	}

	if len(base64Images) == 0 {
		fmt.Println("Không tìm thấy ảnh nào trong folder test_images. Vui lòng thêm ảnh vào đó.")
		return
	}

	fmt.Printf("Tìm thấy %d ảnh. Bắt đầu gửi request...\n", len(base64Images))

	var wg sync.WaitGroup
	start := time.Now()

	results := make(chan string, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		// Lấy xoay vòng ảnh nếu số request > số ảnh
		img := base64Images[i%len(base64Images)]

		go func(id int, base64Img string) {
			defer wg.Done()

			reqBody := VerifyAttendanceRequest{
				Images: base64Images,
			}
			jsonData, _ := json.Marshal(reqBody)

			resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				results <- fmt.Sprintf("Request %d: LỖI KẾT NỐI - %v", id, err)
				return
			}
			defer resp.Body.Close()

			body, _ := ioutil.ReadAll(resp.Body)
			results <- fmt.Sprintf("Request %d: HTTP %d - %s", id, resp.StatusCode, string(body))
		}(i+1, img)
	}

	wg.Wait()
	close(results)

	duration := time.Since(start)

	fmt.Println("\n--- KẾT QUẢ ---")
	for res := range results {
		fmt.Println(res)
	}
	fmt.Printf("\nHoàn thành %d request trong %v\n", numRequests, duration)
}
