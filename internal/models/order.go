package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Confirmed OrderStatus = "confirmed"
	Preparing OrderStatus = "preparing"
	OnTheWay  OrderStatus = "on_the_way"
	Delivered OrderStatus = "delivered"
	Cancelled OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case Pending, Confirmed, Preparing, OnTheWay, Delivered, Cancelled:
		return true
	}
	return false
}

type Order struct {
	ID              uuid.UUID   `gorm:"type:uuid;primary_key;"`
	CustomerID      uuid.UUID   `gorm:"type:uuid;not null"`
	Customer        User        `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	RestaurantID    uuid.UUID   `gorm:"type:uuid;not null"`
	Restaurant      Restaurant  `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	DriverID        *uuid.UUID  `gorm:"type:uuid"`
	Driver          *User       `gorm:"foreignKey:DriverID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	DeliveryAddress string      `gorm:"not null"`
	Total           int         `gorm:"not null"`
	Status          OrderStatus `gorm:"type:varchar(20);default:'pending';not null"`
	Items           []OrderItem `gorm:"foreignKey:OrderID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Quantity  int       `gorm:"not null"`
	UnitPrice int       `gorm:"not null"`
}

type OrderStatusHistory struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key;"`
	OrderID   uuid.UUID   `gorm:"type:uuid;not null"`
	Status    OrderStatus `gorm:"type:varchar(20);not null"`
	ChangedBy uuid.UUID   `gorm:"type:uuid;not null"`
	ChangedAt time.Time
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	o.ID = uuid.New()
	return
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	oi.ID = uuid.New()
	return
}

func (osh *OrderStatusHistory) BeforeCreate(tx *gorm.DB) (err error) {
	osh.ID = uuid.New()
	return
}
