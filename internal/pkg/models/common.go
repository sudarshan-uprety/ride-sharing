// models/common.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Common struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	IsDeleted   bool           `gorm:"default:false"`
	LastLoginAt *time.Time
	CreatedBy   *uint
	UpdatedBy   *uint
	DeletedBy   *uint
}

// BeforeCreate hook to set default UUID if not set
func (c *Common) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
