package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	Customer        UserRole = "customer"
	Driver          UserRole = "driver"
	RestaurantOwner UserRole = "restaurant owner"
)

func (r UserRole) IsValid() bool {
	switch r {
	case Customer, Driver, RestaurantOwner:
		return true
	}
	return false
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      UserRole  `gorm:"default:'customer';type:varchar(20);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
