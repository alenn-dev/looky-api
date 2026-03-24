package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LOGIN | REGISTER
type RegisterDTO struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Role     UserRole `json:"role"`
}

type LoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// USER
type UpdateUserDTO struct {
	Name      *string   `json:"name"`
	Email     *string   `json:"email"`
	Role      *UserRole `json:"role"`
	UpdatedAt time.Time
}

type UserResponse struct {
	ID        uuid.UUID      `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      UserRole       `json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// RESTAURANT
type CreateRestaurantDTO struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type UpdateRestaurantDTO struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Phone    *string `json:"phone"`
	IsActive *bool   `json:"is_active"`
}

type RestaurantResponse struct {
	ID       uuid.UUID `json:"id"`
	OwnerID  uuid.UUID `json:"owner_id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Phone    string    `json:"phone"`
	IsActive bool      `json:"is_active"`
}

// PRODUCT
type CreateProductDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type UpdateProductDTO struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Price       *int    `json:"price"`
	IsAvailable *bool   `json:"is_available"`
}

type ProductResponse struct {
	ID           uuid.UUID `json:"id"`
	RestaurantID uuid.UUID `json:"restaurant_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        int       `json:"price"`
	IsAvailable  bool      `json:"is_available"`
}

// ORDER
type OrderItemDTO struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CreateOrderDTO struct {
	RestaurantID    string         `json:"restaurant_id"`
	DeliveryAddress string         `json:"delivery_address"`
	Items           []OrderItemDTO `json:"items"`
}

type UpdateOrderStatusDTO struct {
	Status OrderStatus `json:"status"`
}

type OrderItemResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	UnitPrice int       `json:"unit_price"`
}

type OrderResponse struct {
	ID              uuid.UUID           `json:"id"`
	CustomerID      uuid.UUID           `json:"customer_id"`
	RestaurantID    uuid.UUID           `json:"restaurant_id"`
	DriverID        *uuid.UUID          `json:"driver_id"`
	DeliveryAddress string              `json:"delivery_address"`
	Total           int                 `json:"total"`
	Status          OrderStatus         `json:"status"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       time.Time           `json:"created_at"`
}

// RESPONSE HISTORY STATUS
type OrderStatusHistoryResponse struct {
	ID        uuid.UUID   `json:"id"`
	OrderID   uuid.UUID   `json:"order_id"`
	Status    OrderStatus `json:"status"`
	ChangedBy uuid.UUID   `json:"changed_by"`
	ChangedAt time.Time   `json:"changed_at"`
}
