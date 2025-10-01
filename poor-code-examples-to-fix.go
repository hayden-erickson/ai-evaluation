package main

import (
	"database/sql"
	"fmt"
	"os"
)

type Config struct {
	APIKey     string
	DBUser     string
	DBPassword string
	DBName     string
}

func loadConfig() (*Config, error) {
	config := &Config{
		APIKey:     os.Getenv("API_KEY"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}

	if config.APIKey == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return config, nil
}

var db *sql.DB

func initDB(config *Config) error {
	var err error
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		config.DBUser, config.DBPassword, config.DBName)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	return nil
}

type User struct {
	Name  string
	Email string
}

// NewUser creates a new User with validation
func NewUser(name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, fmt.Errorf("name and email are required")
	}
	return &User{
		Name:  name,
		Email: email,
	}, nil
}

func (u *User) Save(db *sql.DB, config *Config) error {
	if u.Name == "" || u.Email == "" {
		return fmt.Errorf("user name and email are required")
	}

	_, err := db.Exec(
		"INSERT INTO users (name, email, api_key) VALUES ($1, $2, $3)",
		u.Name, u.Email, config.APIKey,
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}
	return nil
}

func (u *User) Update(db *sql.DB) error {
	if u.Name == "" || u.Email == "" {
		return fmt.Errorf("user name and email are required")
	}

	result, err := db.Exec(
		"UPDATE users SET name = $1 WHERE email = $2",
		u.Name, u.Email,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with email: %s", u.Email)
	}
	return nil
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	var user User
	err := db.QueryRow(
		"SELECT name, email FROM users WHERE id = $1",
		id,
	).Scan(&user.Name, &user.Email)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no user found with id: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	return &user, nil
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := initDB(config); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	user, err := NewUser("John Doe", "john@example.com")
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		os.Exit(1)
	}

	if err := user.Save(db, config); err != nil {
		fmt.Printf("Failed to save user: %v\n", err)
		os.Exit(1)
	}

	updatedUser, err := GetUserByID(db, 1)
	if err != nil {
		fmt.Printf("Failed to get user: %v\n", err)
		os.Exit(1)
	}

	updatedUser.Name = "Jane Doe"
	if err := updatedUser.Update(db); err != nil {
		fmt.Printf("Failed to update user: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully updated user")
}
