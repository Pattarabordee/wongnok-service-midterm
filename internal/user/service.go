package user

import (
	"strings"
	"wongnok/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IService interface {
	UpsertWithClaims(claims model.Claims) (model.User, error)
	GetByID(claims model.Claims) (model.User, error)
	GetRecipes(userID string, claims model.Claims) (model.FoodRecipes, error)
}

type Service struct {
	Repository IRepository
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository: NewRepository(db),
	}
}

func (service Service) UpsertWithClaims(claims model.Claims) (model.User, error) {
	validate := validator.New()
	if err := validate.Struct(claims); err != nil {
		return model.User{}, errors.Wrap(err, "claims invalid")
	}

	user, err := service.Repository.GetByID(claims.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, errors.Wrap(err, "find user")
	}

	// Set claims information to user model
	user = user.FromClaims(claims)

	if err := service.Repository.Upsert(&user); err != nil {
		return model.User{}, errors.Wrap(err, "upsert user")
	}

	return user, nil
}

func (service Service) GetByID(claims model.Claims) (model.User, error) {

	user, err := service.Repository.GetByID(claims.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, errors.Wrap(err, "find user")
	}

	return user, nil
}

func (service Service) GetRecipes(userID string, claims model.Claims) (model.FoodRecipes, error) {

	if strings.ToLower(userID) == "self" {
		userID = claims.ID
	}

	if _, err := service.Repository.GetByID(claims.ID); err != nil {
		return model.FoodRecipes{}, errors.Wrap(err, "find user")
	}

	foodRecipes, err := service.Repository.GetRecipes(userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.FoodRecipes{}, errors.Wrap(err, "get recipes")
	}

	foodRecipes = foodRecipes.CalculateAverageRatings()

	return foodRecipes, nil
}
