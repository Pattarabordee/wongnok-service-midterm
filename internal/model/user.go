package model

import (
	"wongnok/internal/model/dto"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID              string
	FirstName       string
	LastName        string
	NickName        string
	ImageProfileUrl string
}

func (user User) FromClaims(claims Claims) User {
	return User{
		Model:           user.Model,
		ID:              claims.ID,
		FirstName:       claims.FirstName,
		LastName:        claims.LastName,
		NickName:        claims.NickName,
		ImageProfileUrl: claims.ImageProfileUrl,
	}
}

func (user User) ToResponse() dto.UserResponse {
	return dto.UserResponse{
		ID:              user.ID,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		NickName:        user.NickName,
		ImageProfileUrl: user.ImageProfileUrl,
	}
}
