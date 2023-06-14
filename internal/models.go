package internal

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email string
	Pass  string
}

type UserValidation struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreate struct {
	Email string `json:"email"`
	Pass  string `json:"password"`
}

var db *gorm.DB

func init() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})
	fmt.Println("db started")
}

func CreateUser(user UserCreate) {
	db.Create(&User{Email: user.Email, Pass: user.Pass})
}

func UserExistsFromEmail(email string) bool {
	var user User
	db.First(&user, "email = ?", email)
	return user.Email != ""
}

func UserIsValid(email, password string) bool {
	var user User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return false
	}

	return user.Pass == password
}

func GetAllUsers() []User {
	var users []User
	db.Find(&users)
	return users

}
