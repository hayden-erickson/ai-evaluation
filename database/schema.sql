-- MySQL Schema for AI Evaluation Application

-- Create database (run this separately if needed)
-- CREATE DATABASE ai_evaluation CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
-- USE ai_evaluation;

-- Business Users table
CREATE TABLE IF NOT EXISTS business_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    company_uuid VARCHAR(36) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_company_uuid (company_uuid),
    INDEX idx_email (email),
    INDEX idx_deleted_at (deleted_at)
);

-- Sites table
CREATE TABLE IF NOT EXISTS sites (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid VARCHAR(36) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    company_uuid VARCHAR(36) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_uuid (uuid),
    INDEX idx_company_uuid (company_uuid),
    INDEX idx_deleted_at (deleted_at)
);

-- User Sites junction table
CREATE TABLE IF NOT EXISTS user_sites (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    site_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES business_users(id) ON DELETE CASCADE,
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_site (user_id, site_id)
);

-- Units table
CREATE TABLE IF NOT EXISTS units (
    id INT AUTO_INCREMENT PRIMARY KEY,
    site_id INT NOT NULL,
    unit_number VARCHAR(50) NOT NULL,
    rental_state ENUM('available', 'occupied', 'maintenance', 'overlock', 'gatelock', 'prelet') DEFAULT 'available',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE,
    INDEX idx_site_id (site_id),
    INDEX idx_rental_state (rental_state),
    INDEX idx_deleted_at (deleted_at),
    UNIQUE KEY unique_site_unit (site_id, unit_number, deleted_at)
);

-- Gate Access Codes table
CREATE TABLE IF NOT EXISTS gate_access_codes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    access_code VARCHAR(20) NOT NULL,
    unit_id INT NOT NULL,
    user_id INT NOT NULL,
    site_id INT NOT NULL,
    state ENUM('active', 'setup', 'pending', 'inactive', 'removed', 'removing', 'overlocking', 'overlocked', 'remove') DEFAULT 'setup',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES business_users(id) ON DELETE CASCADE,
    FOREIGN KEY (site_id) REFERENCES sites(id) ON DELETE CASCADE,
    INDEX idx_access_code (access_code),
    INDEX idx_unit_id (unit_id),
    INDEX idx_user_id (user_id),
    INDEX idx_site_id (site_id),
    INDEX idx_state (state),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_unit_site (unit_id, site_id)
);

-- Insert sample data for testing
INSERT IGNORE INTO sites (id, uuid, name, company_uuid) VALUES 
(1, '550e8400-e29b-41d4-a716-446655440001', 'Test Site 1', '550e8400-e29b-41d4-a716-446655440000'),
(2, '550e8400-e29b-41d4-a716-446655440002', 'Test Site 2', '550e8400-e29b-41d4-a716-446655440000');

INSERT IGNORE INTO business_users (id, company_uuid, email, first_name, last_name) VALUES 
(1, '550e8400-e29b-41d4-a716-446655440000', 'test@example.com', 'Test', 'User'),
(2, '550e8400-e29b-41d4-a716-446655440000', 'admin@example.com', 'Admin', 'User');

INSERT IGNORE INTO user_sites (user_id, site_id) VALUES 
(1, 1), (1, 2), (2, 1), (2, 2);

INSERT IGNORE INTO units (id, site_id, unit_number, rental_state) VALUES 
(1, 1, 'A101', 'available'),
(2, 1, 'A102', 'occupied'),
(3, 2, 'B201', 'available'),
(4, 2, 'B202', 'maintenance');
