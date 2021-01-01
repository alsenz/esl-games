package model

import (
	"time"

	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

func (base *Base) BeforeCreate(db *gorm.DB) error {
	base.ID = uuid.NewV4()
	return nil
}