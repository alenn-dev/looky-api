package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Restaurant struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null"`
	Owner     User      `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name      string    `gorm:"not null"`
	Address   string    `gorm:"not null"`
	Phone     string    `gorm:"not null"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
