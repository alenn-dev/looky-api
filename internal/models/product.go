package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;not null"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name         string     `gorm:"not null"`
	Description  string
	Price        int  `gorm:"not null"` // In cents (Ex: 15.000,00 colombian pesos or 1500 american cents)
	IsAvailable  bool `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
