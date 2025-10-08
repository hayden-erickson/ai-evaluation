package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var (
	apiKey     string
	dbUser     string
	dbPassword string
	dbName     string
)

var db *sql.DB

func init() {
	// Load configuration from environment variables
	apiKey = os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is required")
	}

	dbUser = os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("DB_USER environment variable is required")
	}

	dbPassword = os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is required")
	}

	dbName = os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("DB_NAME environment variable is required")
	}

	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName))
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

type User struct {
	Name  string
	Email string
	Age   int
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) SetName(name string) {
	u.Name = name
}

func (u *User) SetAge(age int) {
	u.Age = age
}

func (u *User) GetAge() int {
	return u.Age
}

func NewUser(name, email string) *User {
	return &User{Name: name, Email: email}
}

func (u *User) Save(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO users (name, email, api_key) VALUES ($1, $2, $3)", u.Name, u.Email, apiKey)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func (u *User) Update(db *sql.DB) error {
	_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE email = $3", u.Name, u.Email, u.Email)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func GetUser(db *sql.DB, id int) (*User, error) {
	var name, email string
	var apiKeyValue string
	err := db.QueryRow("SELECT name, email, api_key FROM users WHERE id = $1", id).Scan(&name, &email, &apiKeyValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &User{Name: name, Email: email}, nil
}

func main() {
	user := NewUser("John Doe", "john@example.com")
	if err := user.Save(db); err != nil {
		log.Fatalf("Failed to save user: %v", err)
	}

	retrievedUser, err := GetUser(db, 1)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}
	retrievedUser.SetName("Jane Doe")
	if err := retrievedUser.Update(db); err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}
}
