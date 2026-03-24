package services

import (
	"errors"
	"looky/internal/database"
	"looky/internal/kafka"
	"looky/internal/models"

	"github.com/google/uuid"
)

func GetOrders(userID string, role models.UserRole) ([]models.OrderResponse, error) {
	var orders []models.Order

	query := database.DB.Preload("Items").Preload("Items.Product")

	switch role {
	case models.Customer:
		query = query.Where("customer_id = ?", userID)
	case models.Driver:
		query = query.Where("driver_id = ?", userID)
	case models.RestaurantOwner:
		// Search the restaurants from owner and filter by their
		var restaurantIDs []uuid.UUID
		database.DB.Model(&models.Restaurant{}).
			Where("owner_id = ?", userID).
			Pluck("id", &restaurantIDs)
		query = query.Where("restaurant_id IN ?", restaurantIDs)
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, errors.New("couldn't get orders")
	}

	var response []models.OrderResponse
	for _, o := range orders {
		response = append(response, BuildOrderResponse(o))
	}
	return response, nil
}

func GetOrder(id string, userID string, role models.UserRole) (*models.OrderResponse, error) {
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

	response := BuildOrderResponse(order)
	return &response, nil
}

func CreateOrder(customerID string, dto models.CreateOrderDTO) (*models.OrderResponse, error) {
	if len(dto.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}

	parsedRestaurantID, err := uuid.Parse(dto.RestaurantID)
	if err != nil {
		return nil, errors.New("invalid restaurant id")
	}

	// Verificate if restaurant exist and it's active
	var restaurant models.Restaurant
	if err := database.DB.Where("id = ? AND is_active = ?", parsedRestaurantID, true).First(&restaurant).Error; err != nil {
		return nil, errors.New("restaurant not found or inactive")
	}

	// Build the items and calculate total
	var items []models.OrderItem
	total := 0

	for _, itemDTO := range dto.Items {
		parsedProductID, err := uuid.Parse(itemDTO.ProductID)
		if err != nil {
			return nil, errors.New("invalid product id")
		}

		var product models.Product
		if err := database.DB.Where("id = ? AND restaurant_id = ? AND is_available = ?", parsedProductID, parsedRestaurantID, true).First(&product).Error; err != nil {
			return nil, errors.New("product not found or unavailable")
		}

		items = append(items, models.OrderItem{
			ProductID: parsedProductID,
			Quantity:  itemDTO.Quantity,
			UnitPrice: product.Price,
		})
		total += product.Price * itemDTO.Quantity
	}

	order := models.Order{
		CustomerID:      parsedCustomerID,
		RestaurantID:    parsedRestaurantID,
		DeliveryAddress: dto.DeliveryAddress,
		Total:           total,
		Status:          models.Pending,
		Items:           items,
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return nil, errors.New("couldn't create order")
	}

	// Save initial state in the history
	database.DB.Create(&models.OrderStatusHistory{
		OrderID:   order.ID,
		Status:    models.Pending,
		ChangedBy: parsedCustomerID,
		ChangedAt: order.CreatedAt,
	})

	kafka.PublishOrderStatus(kafka.OrderStatusEvent{
		OrderID:    order.ID.String(),
		CustomerID: order.CustomerID.String(),
		Status:     string(models.Pending),
	})

	if err := database.DB.Preload("Items").Preload("Items.Product").Where("id = ?", order.ID).First(&order).Error; err != nil {
		return nil, errors.New("couldn't find order")
	}

	response := BuildOrderResponse(order)
	return &response, nil
}

func UpdateOrderStatus(id string, userID string, role models.UserRole, dto models.UpdateOrderStatusDTO) (*models.OrderResponse, error) {
	if !dto.Status.IsValid() {
		return nil, errors.New("invalid status")
	}

	var order models.Order
	if err := database.DB.Where("id = ?", id).First(&order).Error; err != nil {
		return nil, errors.New("order not found")
	}

	// Validate who can change to which state
	if err := validateStatusTransition(order.Status, dto.Status, role); err != nil {
		return nil, err
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// If the delivery person accepts the order, it is automatically assigned
	if dto.Status == models.OnTheWay && role == models.Driver {
		order.DriverID = &parsedUserID
	}

	order.Status = dto.Status
	if err := database.DB.Save(&order).Error; err != nil {
		return nil, errors.New("couldn't update order status")
	}

	// Save history
	database.DB.Create(&models.OrderStatusHistory{
		OrderID:   order.ID,
		Status:    dto.Status,
		ChangedBy: parsedUserID,
	})

	kafka.PublishOrderStatus(kafka.OrderStatusEvent{
		OrderID:    order.ID.String(),
		CustomerID: order.CustomerID.String(),
		Status:     string(dto.Status),
	})

	if err := database.DB.Preload("Items").Preload("Items.Product").Where("id = ?", order.ID).First(&order).Error; err != nil {
		return nil, errors.New("couldn't find order")
	}

	response := BuildOrderResponse(order)
	return &response, nil
}

// validateStatusTransition defines who can change to which state
func validateStatusTransition(current models.OrderStatus, next models.OrderStatus, role models.UserRole) error {
	allowed := map[models.UserRole]map[models.OrderStatus][]models.OrderStatus{
		models.RestaurantOwner: {
			models.Pending:   {models.Confirmed, models.Cancelled},
			models.Confirmed: {models.Preparing, models.Cancelled},
			models.Preparing: {models.OnTheWay},
		},
		models.Driver: {
			models.Preparing: {models.OnTheWay},
			models.OnTheWay:  {models.Delivered},
		},
		models.Customer: {
			models.Pending: {models.Cancelled},
		},
	}

	transitions, ok := allowed[role]
	if !ok {
		return errors.New("unauthorized to change order status")
	}

	validNext, ok := transitions[current]
	if !ok {
		return errors.New("cannot change status from " + string(current))
	}

	for _, s := range validNext {
		if s == next {
			return nil
		}
	}

	return errors.New("invalid status transition")
}

// BuildOrderResponse builds the response with the items
func BuildOrderResponse(o models.Order) models.OrderResponse {
	var items []models.OrderItemResponse
	for _, item := range o.Items {
		items = append(items, models.OrderItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Product.Name,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	return models.OrderResponse{
		ID:              o.ID,
		CustomerID:      o.CustomerID,
		RestaurantID:    o.RestaurantID,
		DriverID:        o.DriverID,
		DeliveryAddress: o.DeliveryAddress,
		Total:           o.Total,
		Status:          o.Status,
		Items:           items,
		CreatedAt:       o.CreatedAt,
	}
}
