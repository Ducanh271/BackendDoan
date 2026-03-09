# Makefile cho dự án Go

.PHONY: help run dev build tidy test clean

help:               ## Hiển thị các lệnh có sẵn
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

run:                ## Chạy ứng dụng (không reload, ổn định)
	go run ./cmd/main.go

dev:                ## Chạy dev mode với auto-reload (dùng air - khuyến nghị)
	air

build:              ## Build binary (tên mặc định là tên folder hoặc chỉ định)
	go build -o bin/app ./cmd/main.go

tidy:               ## Tidy modules và format code
	go mod tidy
	go fmt ./...

test:               ## Chạy tất cả test
	go test -v ./...

clean:              ## Xóa binary, cache...
	rm -rf bin/
	go clean -cache -testcache

# Nếu main.go nằm ở cmd/server hoặc cmd/api thì thay ./cmd/main.go thành ./cmd/server hoặc ./cmd/api
