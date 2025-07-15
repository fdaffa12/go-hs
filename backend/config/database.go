package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// GetDatabaseConfig returns database configuration from environment variables
func GetDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "cors_auth_db"),
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// InitDB initializes database connection
func InitDB() (*sql.DB, error) {
	config := GetDatabaseConfig()

	// First, connect without database to create it if it doesn't exist
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)

	tempDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}
	defer tempDB.Close()

	// Create database if it doesn't exist
	_, err = tempDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}

	// Now connect to the specific database
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	// Create tables if they don't exist
	if err = createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

// createTables creates necessary tables
func createTables() error {
	userTableSQL := `
	CREATE TABLE IF NOT EXISTS hs_wsb_user (
		NIK VARCHAR(10) PRIMARY KEY,
		NAME VARCHAR(25) NOT NULL,
		EMAIL VARCHAR(25) NOT NULL UNIQUE,
		PASSWORD VARCHAR(255) NOT NULL,
		LEVEL INT NOT NULL DEFAULT 1,
		PROFILE_PICTURE VARCHAR(500) NULL,
		remember_token VARCHAR(100) NULL,
		created_at TIMESTAMP NULL,
		updated_at TIMESTAMP NULL
	);
	`

	_, err := db.Exec(userTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	// Add departments table
	departmentTableSQL := `
	CREATE TABLE IF NOT EXISTS departments (
		DEPT_SHORT_NAME VARCHAR(5) PRIMARY KEY,
		DEPT_LONG_NAME VARCHAR(50) NOT NULL,
		DELETE_STATUS BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(departmentTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create departments table: %v", err)
	}

	// Add employees table
	employeeTableSQL := `
	CREATE TABLE IF NOT EXISTS employees (
		NIK VARCHAR(10) PRIMARY KEY,
		NAME VARCHAR(22) NOT NULL,
		DEPT_SHORT_NAME VARCHAR(5) NOT NULL,
		ENTERANCE_DATE DATE NOT NULL,
		TITLE VARCHAR(22) NOT NULL,
		GENDER TINYINT(1) NOT NULL COMMENT '1=Laki-laki, 0=Perempuan',
		RFID_ID VARCHAR(15),
		WORKING_AREA TINYINT(1) NOT NULL DEFAULT 0 COMMENT '1=ALL AREA, 2=F1, 3=F2, 4=F3',
		DELETE_STATUS TINYINT(1) NOT NULL DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (DEPT_SHORT_NAME) REFERENCES departments(DEPT_SHORT_NAME)
	);
	`

	_, err = db.Exec(employeeTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create employees table: %v", err)
	}

	// Add profile_picture column if it doesn't exist (for existing databases)
	// Check if column exists first
	var columnExists int
	checkQuery := `SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'profile_picture'`
	if err := db.QueryRow(checkQuery).Scan(&columnExists); err != nil {
		fmt.Printf("Error checking if profile_picture column exists: %v\n", err)
	} else if columnExists == 0 {
		alterQuery := `ALTER TABLE users ADD COLUMN profile_picture VARCHAR(500) NULL`
		if _, err := db.Exec(alterQuery); err != nil {
			fmt.Printf("Error adding profile_picture column: %v\n", err)
		} else {
			fmt.Println("profile_picture column added successfully")
		}
	}

	return nil
}

// GetDB returns database instance
func GetDB() *sql.DB {
	return db
}

// CloseDB closes database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}