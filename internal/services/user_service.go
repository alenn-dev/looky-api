package services

import (
	"errors"
	"looky/internal/database"
	"looky/internal/models"
	"time"
)

func GetMe(userID string) (*models.UserResponse, error) {
	var user models.UserResponse
	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("couldn't find user")
	}

	return &user, nil
}

func UpdateUser(userID string, dto models.UpdateUserDTO) (*models.UserResponse, error) {
	if dto.Role != nil && !dto.Role.IsValid() {
		return nil, errors.New("invalid role")
	}

	dto.UpdatedAt = time.Now()
	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(&dto)
	if result.Error != nil {
		return nil, errors.New("couldn't update user")
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	var user models.UserResponse
	if err := database.DB.Model(&models.User{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, errors.New("couldn't find user")
	}

	return &user, nil
}
