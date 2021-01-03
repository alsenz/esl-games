package model

import (
	"time"

	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time	 `json:"createdAt"`
	UpdatedAt time.Time	 `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

func (base *Base) BeforeCreate(db *gorm.DB) error {
	base.ID = uuid.NewV4()
	return nil
}