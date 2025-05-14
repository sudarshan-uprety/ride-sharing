package models

import (
	CommonModels "ride-sharing/internal/pkg/models" // Import the common model package
	"time"
)

type Rider struct {
	CommonModels.Common
	FullName          string    `gorm:"not null"`
	Phone             string    `gorm:"unique;not null"`
	Email             string    `gorm:"unique;not null"`
	Password          string    `gorm:"not null"`
	LicenseNumber     string    `gorm:"unique;not null"`
	LicenseIssueDate  time.Time `gorm:"not null"`
	LicenseExpiryDate time.Time `gorm:"not null"`
	LicenseCategory   string    `gorm:"type:enum('A', 'B', 'K');not null"` // A: Bike, B: Car, K: Scooter
	BlueBookNumber    string    `gorm:"unique;not null"`                   // Vehicle registration
	VehicleType       string    `gorm:"not null"`                          // bike, car, premium, xl
	VehicleModel      string    `gorm:"type:enum('A', 'B', 'K');not null"` // A: Bike, B: Car, K: Scooter
	VehicleYear       int       `gorm:"not null"`
	IsApproved        bool      `gorm:"default:false"`
	ApprovalStatus    string    `gorm:"type:enum('pending', 'approved', 'rejected');default:'pending'"`
	Rating            float64   `gorm:"default:0.0"`
	TotalTrips        int       `gorm:"default:0"`
	OnlineStatus      bool      `gorm:"default:false"`
	PasswordChangedAt *time.Time
}
