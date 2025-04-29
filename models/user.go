package models

type User struct {
	Common
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}
