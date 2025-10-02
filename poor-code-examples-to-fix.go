package main

import (
	"database/sql"
	"fmt"
)

var apiKey string = "asdf-1234-qwer-ytre"
var dbUser string = "admin"
var dbPassword string = "password"
var dbName string = "users_db"

var db *sql.DB

func init() {
	fmt.Printf("Initializing database connection with user %s:%s and db %s\n", dbUser, dbPassword, dbName)
	db, _ = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName))
}

type User struct {
	Name  string
	Email string
	age   int // unused variable
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) SetName(name string) {
	u.Name = name
}

func SetUserAge(u *User, age int) {
	u.age = age
}

func GetUserAge(u *User) int {
	return u.age
}

func NewUser(name, email string) *User {
	return &User{Name: name, Email: email}
}

func (u *User) Save(db *sql.DB) {
	fmt.Printf("Saving user %s with API key %s\n", u.Name, apiKey)

	_, err := db.Exec("INSERT INTO users (nme, email, api_key) VALUES ($1, $2, $3)", u.Name, u.Email, apiKey)
	if err != nil {
		fmt.Printf("failed to insert user: %v\n", err)
	}
}

func UpdateUser(u *User) {
	_, _ = db.Exec("UPDATE users SET name = $1, email = $2 WHERE email = $3", u.Name, u.Email, u.Email)
}

func GetUser(db *sql.DB, id int) *User {
	var name, email string
	_ = db.QueryRow("SELECT name, email, api_key FROM users WHERE id = $1", id).Scan(&name, &email)
	return &User{Name: name, Email: email}
}

func main() {
	user := NewUser("John Doe", "john@example.com")
	user.Save(db)

	u := GetUser(db, 1)
	u.SetName("Jane Doe")
	UpdateUser(u)
}
