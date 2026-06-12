-- database/schema.sql

CREATE DATABASE IF NOT EXISTS AttendanceDB;
USE AttendanceDB;

-- 1. Bảng Chức vụ trước)
CREATE TABLE positions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    salary_grade VARCHAR(50)
);

-- 2. Bảng Phòng ban 
CREATE TABLE departments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    manager_id INT NULL
);

-- 3. Bảng Nhân viên
CREATE TABLE employees (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    department_id INT,
    position_id INT,
    hire_date DATE,
    embedding JSON, -- Lưu vector khuôn mặt từ Python (128 chiều)
    status VARCHAR(50) DEFAULT 'ACTIVE', -- ACTIVE, INACTIVE, SUSPENDED

    -- Thông tin xác thực / đăng nhập
    password_hash VARCHAR(255),          -- Bcrypt hash của mật khẩu đăng nhập
    role VARCHAR(50) DEFAULT 'user',     -- user, ADMIN
    is_first_login BOOLEAN DEFAULT TRUE, -- TRUE: bắt buộc đổi mật khẩu mặc định qua OTP trước khi dùng
    otp_code VARCHAR(10),                -- Mã OTP đổi mật khẩu lần đầu
    otp_expires_at DATETIME,             -- Thời điểm hết hạn của OTP

    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
    FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL
);

-- GÁN KHÓA NGOẠI CHO BẢNG DEPARTMENTS:
ALTER TABLE departments 
ADD CONSTRAINT fk_department_manager 
FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL;

-- 4. Bảng Nhật ký điểm danh
CREATE TABLE attendance_records (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id INT NOT NULL,
    date DATE NOT NULL,
    check_in DATETIME NOT NULL,
    check_out DATETIME NULL,
    location VARCHAR(255),
    method VARCHAR(50) DEFAULT 'FACE_RECOGNITION', -- FACE_RECOGNITION, MANUAL, PASSWORD
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

-- 5. Bảng Refresh Token (quản lý phiên đăng nhập)
CREATE TABLE refresh_tokens (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id INT NOT NULL,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE
);

-- 6. Bảng Nhật ký truy cập (log mọi lần điểm danh thành công/thất bại)
CREATE TABLE access_logs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    employee_id INT NULL, -- NULL khi không nhận diện được người (truy cập lạ)
    access_point VARCHAR(255),            -- VD: "Mobile App", "Tablet cửa chính"
    action_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50),                   -- GRANTED, DENIED
    confidence_distance DOUBLE,           -- Khoảng cách Euclidean khi so khớp khuôn mặt
    failure_reason VARCHAR(255),
    FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE SET NULL
);

-- 7. Bảng Cấu hình công ty (tọa độ trụ sở + bán kính cho phép điểm danh)
CREATE TABLE company_configs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    latitude DOUBLE NOT NULL DEFAULT 0,
    longitude DOUBLE NOT NULL DEFAULT 0,
    max_radius DOUBLE NOT NULL DEFAULT 0, -- Bán kính cho phép (mét)
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 8. Bảng Wi-Fi hợp lệ của công ty (đối chiếu BSSID khi điểm danh qua mobile)
CREATE TABLE company_wifis (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_config_id INT NOT NULL,
    bssid VARCHAR(100) NOT NULL UNIQUE,
    description VARCHAR(255),
    FOREIGN KEY (company_config_id) REFERENCES company_configs(id) ON DELETE CASCADE
);
