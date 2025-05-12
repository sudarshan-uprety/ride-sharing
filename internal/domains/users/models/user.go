package models

import (
	CommonModels "ride-sharing/internal/pkg/models" // Import the common model package
	"time"
)

type User struct {
	CommonModels.Common `swaggerignore:"true"`
	FullName            string `gorm:"not null"`
	Phone               string `gorm:"unique;not null"`
	Address             string `gorm:"not null"`
	Email               string `gorm:"unique;not null"`
	Password            string `gorm:"not null"`
	Active              bool   `gorm:"default:false"`
	PasswordChangedAt   *time.Time
}

func (User) TableName() string {
	return "users"
}
