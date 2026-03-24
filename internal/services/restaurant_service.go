package services

import (
	"errors"
	"looky/internal/database"
	"looky/internal/models"

	"github.com/google/uuid"
)

func GetRestaurants() ([]models.RestaurantResponse, error) {
	var restaurants []models.RestaurantResponse
	if err := database.DB.Model(&models.Restaurant{}).Where("is_active = ?", true).Find(&restaurants).Error; err != nil {
		return nil, errors.New("couldn't get restaurants")
	}
	return restaurants, nil
}

func GetRestaurant(id string) (*models.RestaurantResponse, error) {
	var restaurant models.RestaurantResponse
	if err := database.DB.Model(&models.Restaurant{}).Where("id = ?", id).First(&restaurant).Error; err != nil {
		return nil, errors.New("restaurant not found")
	}
	return &restaurant, nil
}

func CreateRestaurant(ownerID string, dto models.CreateRestaurantDTO) (*models.RestaurantResponse, error) {
	parsedID, err := uuid.Parse(ownerID)
	if err != nil {
		return nil, errors.New("invalid owner id")
	}

	restaurant := models.Restaurant{
		OwnerID: parsedID,
		Name:    dto.Name,
		Address: dto.Address,
		Phone:   dto.Phone,
	}

	if err := database.DB.Create(&restaurant).Error; err != nil {
		return nil, errors.New("couldn't create restaurant")
	}

	var response models.RestaurantResponse
	if err := database.DB.Model(&models.Restaurant{}).Where("id = ?", restaurant.ID).First(&response).Error; err != nil {
		return nil, errors.New("couldn't find restaurant")
	}

	return &response, nil
}

func UpdateRestaurant(id, ownerID string, dto models.UpdateRestaurantDTO) (*models.RestaurantResponse, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid restaurant id")
	}

	var restaurant models.Restaurant
	if err := database.DB.Where("id = ? AND owner_id = ?", parsedID, ownerID).First(&restaurant).Error; err != nil {
		return nil, errors.New("restaurant not found or unathorized")
	}

	if err := database.DB.Model(&models.Restaurant{ID: parsedID}).Updates(&dto).Error; err != nil {
		return nil, errors.New("couldn't update restaurant")
	}

	var response models.RestaurantResponse
	if err := database.DB.Model(&models.Restaurant{}).Where("id = ?", parsedID).First(&response).Error; err != nil {
		return nil, errors.New("couldn't find restaurant")
	}

	return &response, nil
}

func DeleteRestaurant(id string, ownerID string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid restaurant id")
	}

	result := database.DB.Where("id = ? AND owner_id = ?", parsedID, ownerID).Delete(&models.Restaurant{})
	if result.Error != nil {
		return errors.New("couldn't delete restaurant")
	}
	if result.RowsAffected == 0 {
		return errors.New("restaurant not found or unauthorized")
	}

	return nil
}
