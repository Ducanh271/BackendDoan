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
