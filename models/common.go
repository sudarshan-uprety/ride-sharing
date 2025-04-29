package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Common struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt         time.Time      `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updatedAt" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
	IsArchived        bool           `json:"isArchived" gorm:"default:false"`
	PasswordChangedAt *time.Time     `json:"passwordChangedAt,omitempty"`
	LastLoginAt       *time.Time     `json:"lastLoginAt,omitempty"`
	CreatedBy         *uint          `json:"createdBy,omitempty"`
	UpdatedBy         *uint          `json:"updatedBy,omitempty"`
	DeletedBy         *uint          `json:"deletedBy,omitempty"`
}
