package services

import (
	"errors"
	"looky/internal/database"
	"looky/internal/models"

	"github.com/google/uuid"
)

func GetProducts(restaurantID string) ([]models.ProductResponse, error) {
	var products []models.ProductResponse
	if err := database.DB.Model(&models.Product{}).
		Where("restaurant_id = ? AND is_available = ?", restaurantID, true).
		Find(&products).Error; err != nil {
		return nil, errors.New("couldn't get products")
	}
	return products, nil
}

func GetProduct(id string, restaurantID string) (*models.ProductResponse, error) {
	var product models.ProductResponse
	if err := database.DB.Model(&models.Product{}).
		Where("id = ? AND restaurant_id = ?", id, restaurantID).
		First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}
	return &product, nil
}

func CreateProduct(restaurantID string, ownerID string, dto models.CreateProductDTO) (*models.ProductResponse, error) {
	var restaurant models.Restaurant
	if err := database.DB.Where("id = ? AND owner_id = ?", restaurantID, ownerID).First(&restaurant).Error; err != nil {
		return nil, errors.New("restaurant not found or unauthorized")
	}

	parsedRestaurantID, err := uuid.Parse(restaurantID)
	if err != nil {
		return nil, errors.New("invalid restaurant id")
	}

	product := models.Product{
		RestaurantID: parsedRestaurantID,
		Name:         dto.Name,
		Description:  dto.Description,
		Price:        dto.Price,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		return nil, errors.New("couldn't create product")
	}

	var response models.ProductResponse
	if err := database.DB.Model(&models.Product{}).Where("id = ?", product.ID).First(&response).Error; err != nil {
		return nil, errors.New("couldn't find product")
	}

	return &response, nil
}

func UpdateProduct(id string, restaurantID string, ownerID string, dto models.UpdateProductDTO) (*models.ProductResponse, error) {
	var restaurant models.Restaurant
	if err := database.DB.Where("id = ? AND owner_id = ?", restaurantID, ownerID).First(&restaurant).Error; err != nil {
		return nil, errors.New("restaurant not found or unauthorized")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid product id")
	}

	var product models.Product
	if err := database.DB.Where("id = ? AND restaurant_id = ?", parsedID, restaurantID).First(&product).Error; err != nil {
		return nil, errors.New("product not found or unauthorized")
	}

	if err := database.DB.Model(&models.Product{ID: parsedID}).Updates(&dto).Error; err != nil {
		return nil, errors.New("couldn't update product")
	}

	var response models.ProductResponse
	if err := database.DB.Model(&models.Product{}).Where("id = ?", parsedID).First(&response).Error; err != nil {
		return nil, errors.New("couldn't find product")
	}

	return &response, nil
}

func DeleteProduct(id string, restaurantID string, ownerID string) error {
	var restaurant models.Restaurant
	if err := database.DB.Where("id = ? AND owner_id = ?", restaurantID, ownerID).First(&restaurant).Error; err != nil {
		return errors.New("restaurant not found or unauthorized")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid product id")
	}

	result := database.DB.Where("id = ? AND restaurant_id = ?", parsedID, restaurantID).Delete(&models.Product{})
	if result.Error != nil {
		return errors.New("couldn't delete product")
	}
	if result.RowsAffected == 0 {
		return errors.New("product not found or unauthorized")
	}

	return nil
}
