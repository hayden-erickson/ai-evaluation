package main

import (
	"database/sql"
	"fmt"
	"os"
)

var apiKey = getEnv("API_KEY", "")
var dbUser = getEnv("DB_USER", "")
var dbPassword = getEnv("DB_PASSWORD", "")
var dbName = getEnv("DB_NAME", "")

var db *sql.DB

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func init() {
	if dbUser == "" || dbPassword == "" || dbName == "" {
		fmt.Println("Warning: Database credentials not set in environment variables")
		return
	}
	fmt.Printf("Initializing database connection with user %s and db %s\n", dbUser, dbName)
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName))
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
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

func (u *User) GetAge() int {
	return u.Age
}

func (u *User) SetAge(age int) {
	u.Age = age
}

func NewUser(name, email string) *User {
	return &User{Name: name, Email: email}
}

func (u *User) Save(db *sql.DB) error {
	fmt.Printf("Saving user %s\n", u.Name)

	_, err := db.Exec("INSERT INTO users (name, email, api_key) VALUES ($1, $2, $3)", u.Name, u.Email, apiKey)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

func UpdateUser(u *User) error {
	_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE email = $3", u.Name, u.Email, u.Email)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func GetUser(db *sql.DB, id int) (*User, error) {
	var name, email string
	err := db.QueryRow("SELECT name, email FROM users WHERE id = $1", id).Scan(&name, &email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &User{Name: name, Email: email}, nil
}

func main() {
	user := NewUser("John Doe", "john@example.com")
	if err := user.Save(db); err != nil {
		fmt.Printf("Error saving user: %v\n", err)
		return
	}

	u, err := GetUser(db, 1)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err)
		return
	}
	u.SetName("Jane Doe")
	if err := UpdateUser(u); err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return
	}
}
