package rating

import (
	"wongnok/internal/model"
	"wongnok/internal/model/dto"
	"wongnok/internal/user"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IUserService user.IService

type IService interface {
	Get(recipeID int) (model.Ratings, error)

	Create(request dto.RatingRequest, recipeID int, claims model.Claims) (model.Rating, error)
}

type Service struct {
	Repository  IRepository
	UserService IUserService
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository:  NewRepository(db),
		UserService: user.NewService(db),
	}
}

func (service Service) Get(recipeID int) (model.Ratings, error) {
	ratings, err := service.Repository.Get(recipeID)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func (service Service) Create(request dto.RatingRequest, recipeID int, claims model.Claims) (model.Rating, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.Rating{}, errors.Wrap(err, "request invalid")
	}

	userID, err := service.UserService.GetByID(claims)
	if err != nil {
		return model.Rating{}, errors.Wrap(err, "create rating")
	}

	var rating model.Rating
	rating = rating.FromRequest(request)
	rating.FoodRecipeID = uint(recipeID)

	rating.UserID = userID.ID

	if err := service.Repository.Create(&rating); err != nil {
		return model.Rating{}, errors.Wrap(err, "create rating")
	}

	return rating, nil
}
