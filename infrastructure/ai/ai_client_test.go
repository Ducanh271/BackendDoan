// ai_client_test.go
package ai

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestAIClient_ConcurrencyLimit(t *testing.T) {
	maxConcurrent := 2 // Giới hạn 2 luồng
	totalRequests := 5 // Tổng cộng 5 người gửi
	var activeRequests int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentActive := atomic.AddInt32(&activeRequests, 1)

		defer atomic.AddInt32(&activeRequests, -1)

		if currentActive > int32(maxConcurrent) {
			t.Errorf("LỖI BẢO MẬT: Có %d luồng lọt vào AI cùng lúc (Cho phép tối đa: %d)", currentActive, maxConcurrent)
		}

		time.Sleep(1 * time.Second) // Giả lập AI xử lý mất 1 giây

		// Trả về JSON giả
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer mockServer.Close()

	// 2. Khởi tạo AI Client trỏ vào Server Ảo
	client := NewAIClient(mockServer.URL, maxConcurrent)
	var wg sync.WaitGroup

	start := time.Now()

	// 3. Bắn cùng lúc 5 request
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = client.VerifyFace([][]byte{[]byte("fake_image")})
		}()
	}

	wg.Wait() // Chờ cả 5 thằng chạy xong
	duration := time.Since(start)

	// Phân tích: Có 5 request, làm 2 cái một lượt (mất 1s) -> Cần 3 lượt (lượt 1: 2 req, lượt 2: 2 req, lượt 3: 1 req)
	// Vậy tổng thời gian chắc chắn phải >= 3 giây. Nếu nó chạy xong trong 1 giây nghĩa là Semaphore bị hỏng!
	if duration < 3*time.Second {
		t.Errorf("LỖI: Semaphore không chặn request! Thời gian chạy quá nhanh: %v", duration)
	} else {
		t.Logf("THÀNH CÔNG: Semaphore hoạt động chuẩn. Thời gian xử lý 5 request: %v", duration)
	}
}
