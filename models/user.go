package models

type User struct {
	Common
	FullName string `json:"full_name"`
	Phone    string `json:"phone" grom:"unique"`
	Address  string `json:"address" grom:"unique"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
}
