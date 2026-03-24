package graphql

import (
	"context"
	"errors"
	"looky/internal/database"
	"looky/internal/models"
	"looky/internal/services"
	"time"

	"github.com/google/uuid"
)

type Resolver struct{}

func (r *Resolver) Query() QueryResolver         { return &queryResolver{r} }
func (r *Resolver) Order() OrderResolver         { return &orderResolver{r} }
func (r *Resolver) OrderItem() OrderItemResolver { return &orderItemResolver{r} }
func (r *Resolver) OrderStatusHistory() OrderStatusHistoryResolver {
	return &orderStatusHistoryResolver{r}
}

type orderResolver struct{ *Resolver }
type orderItemResolver struct{ *Resolver }
type orderStatusHistoryResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

func (r *queryResolver) Orders(ctx context.Context, status *models.OrderStatus, from *string, to *string) ([]*models.OrderResponse, error) {
	userID, role, err := extractClaims(ctx)
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	query := database.DB.Preload("Items").Preload("Items.Product")

	switch role {
	case models.Customer:
		query = query.Where("customer_id = ?", userID)
	case models.Driver:
		query = query.Where("driver_id = ?", userID)
	case models.RestaurantOwner:
		var restaurantIDs []uuid.UUID
		database.DB.Model(&models.Restaurant{}).
			Where("owner_id = ?", userID).
			Pluck("id", &restaurantIDs)
		query = query.Where("restaurant_id IN ?", restaurantIDs)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if from != nil {
		query = query.Where("created_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("created_at <= ?", *to)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, errors.New("couldn't get orders")
	}

	var response []*models.OrderResponse
	for _, o := range orders {
		r := services.BuildOrderResponse(o)
		response = append(response, &r)
	}
	return response, nil
}

func (r *queryResolver) Order(ctx context.Context, id string) (*models.OrderResponse, error) {
	userID, role, err := extractClaims(ctx)
	if err != nil {
		return nil, err
	}

	var order models.Order
	query := database.DB.Preload("Items").Preload("Items.Product").Where("id = ?", id)

	switch role {
	case models.Customer:
		query = query.Where("customer_id = ?", userID)
	case models.Driver:
		query = query.Where("driver_id = ?", userID)
	}

	if err := query.First(&order).Error; err != nil {
		return nil, errors.New("order not found or unauthorized")
	}

	response := services.BuildOrderResponse(order)
	return &response, nil
}

func (r *queryResolver) OrderHistory(ctx context.Context, orderID string) ([]*models.OrderStatusHistoryResponse, error) {
	_, _, err := extractClaims(ctx)
	if err != nil {
		return nil, err
	}

	var history []models.OrderStatusHistory
	if err := database.DB.Where("order_id = ?", orderID).
		Order("changed_at ASC").
		Find(&history).Error; err != nil {
		return nil, errors.New("couldn't get order history")
	}

	var response []*models.OrderStatusHistoryResponse
	for _, h := range history {
		item := models.OrderStatusHistoryResponse{
			ID:        h.ID,
			OrderID:   h.OrderID,
			Status:    h.Status,
			ChangedBy: h.ChangedBy,
			ChangedAt: h.ChangedAt,
		}
		response = append(response, &item)
	}
	return response, nil
}

// Helpers
// OrderResolver — convierte uuid.UUID a string
func (r *orderResolver) ID(ctx context.Context, obj *models.OrderResponse) (string, error) {
	return obj.ID.String(), nil
}
func (r *orderResolver) CustomerID(ctx context.Context, obj *models.OrderResponse) (string, error) {
	return obj.CustomerID.String(), nil
}
func (r *orderResolver) RestaurantID(ctx context.Context, obj *models.OrderResponse) (string, error) {
	return obj.RestaurantID.String(), nil
}
func (r *orderResolver) DriverID(ctx context.Context, obj *models.OrderResponse) (*string, error) {
	if obj.DriverID == nil {
		return nil, nil
	}
	s := obj.DriverID.String()
	return &s, nil
}
func (r *orderResolver) CreatedAt(ctx context.Context, obj *models.OrderResponse) (string, error) {
	return obj.CreatedAt.Format(time.RFC3339), nil
}

// OrderItemResolver
func (r *orderItemResolver) ID(ctx context.Context, obj *models.OrderItemResponse) (string, error) {
	return obj.ID.String(), nil
}
func (r *orderItemResolver) ProductID(ctx context.Context, obj *models.OrderItemResponse) (string, error) {
	return obj.ProductID.String(), nil
}

// OrderStatusHistoryResolver
func (r *orderStatusHistoryResolver) ID(ctx context.Context, obj *models.OrderStatusHistoryResponse) (string, error) {
	return obj.ID.String(), nil
}
func (r *orderStatusHistoryResolver) OrderID(ctx context.Context, obj *models.OrderStatusHistoryResponse) (string, error) {
	return obj.OrderID.String(), nil
}
func (r *orderStatusHistoryResolver) ChangedBy(ctx context.Context, obj *models.OrderStatusHistoryResponse) (string, error) {
	return obj.ChangedBy.String(), nil
}
func (r *orderStatusHistoryResolver) ChangedAt(ctx context.Context, obj *models.OrderStatusHistoryResponse) (string, error) {
	return obj.ChangedAt.Format(time.RFC3339), nil
}
