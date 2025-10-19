package foodrecipe

import (
	"wongnok/internal/global"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IService interface {
	Create(request dto.FoodRecipeRequest, claims model.Claims) (model.FoodRecipe, error)
	Get(foodRecipeQuery model.FoodRecipeQuery) (model.FoodRecipes, int64, error)
	GetByID(id int) (model.FoodRecipe, error)
	Update(request dto.FoodRecipeRequest, id int, claims model.Claims) (model.FoodRecipe, error)
	Delete(id int, claims model.Claims) error
}

type Service struct {
	Repository IRepository
}

func NewService(db *gorm.DB) IService {
	return &Service{
		Repository: NewRepository(db),
	}
}

func (service Service) Create(request dto.FoodRecipeRequest, claims model.Claims) (model.FoodRecipe, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "request invalid")
	}

	var recipe model.FoodRecipe
	recipe = recipe.FromRequest(request, claims)

	if err := service.Repository.Create(&recipe); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "create recipe")
	}

	return recipe, nil
}

func (service Service) Get(foodRecipeQuery model.FoodRecipeQuery) (model.FoodRecipes, int64, error) {
	total, err := service.Repository.Count()
	if err != nil {
		return nil, 0, err
	}

	results, err := service.Repository.Get(foodRecipeQuery)
	if err != nil {
		return nil, 0, err
	}

	results = results.CalculateAverageRatings()

	return results, total, nil
}

func (service Service) GetByID(id int) (model.FoodRecipe, error) {
	results, err := service.Repository.GetByID(id)
	if err != nil {
		return model.FoodRecipe{}, err
	}

	results = results.CalculateAverageRating()

	return results, nil
}

func (service Service) Update(request dto.FoodRecipeRequest, id int, claims model.Claims) (model.FoodRecipe, error) {
	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "request invalid")
	}

	recipe, err := service.Repository.GetByID(id)
	if err != nil {
		// กรณีไม่พบ id ที่ต้องการ update
		return model.FoodRecipe{}, errors.Wrap(err, "find recipe")
	}

	if recipe.UserID != claims.ID {
		// กรณี user ที่ login ไม่ตรงกับ user ที่สร้าง recipe
		return model.FoodRecipe{}, global.ErrForbidden
	}

	recipe = recipe.FromRequest(request, claims)

	if err := service.Repository.Update(&recipe); err != nil {
		return model.FoodRecipe{}, errors.Wrap(err, "update recipe")
	}

	recipe = recipe.CalculateAverageRating()

	return recipe, nil
}

func (service Service) Delete(id int, claims model.Claims) error {
	recipe, err := service.Repository.GetByID(id)
	if err != nil {
		// กรณีไม่พบ id ที่ต้องการ update
		return errors.Wrap(err, "find recipe")

	}

	if recipe.UserID != claims.ID {
		// กรณี user ที่ login ไม่ตรงกับ user ที่สร้าง recipe
		return global.ErrForbidden
	}

	return service.Repository.Delete(id)
}
