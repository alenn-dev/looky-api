package services

import (
	"errors"
	"looky/internal/database"
	"looky/internal/models"
	"looky/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func Register(dto models.RegisterDTO) (*models.User, error) {
	existing := models.User{}
	if err := database.DB.Where("email = ?", dto.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already in use")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("couldn't process password")
	}

	user := models.User{
		Name:     dto.Name,
		Email:    dto.Email,
		Password: string(hashed),
		Role:     dto.Role,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, errors.New("couldn't create user")
	}

	user.Password = ""
	return &user, nil
}

func Login(dto models.LoginDTO) (string, error) {
	user := models.User{}
	if err := database.DB.Where("email = ?", dto.Email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateJWT(user.ID.String(), string(user.Role))
	if err != nil {
		return "", errors.New("couldn't generate token")
	}

	return token, nil
}
