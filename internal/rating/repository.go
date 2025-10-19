package rating

import (
	"wongnok/internal/model"

	"gorm.io/gorm"
)

type IRepository interface {
	Get(recipeID int) (model.Ratings, error)
	Create(rating *model.Rating) error
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		DB: db,
	}
}

func (repo Repository) Get(recipeID int) (model.Ratings, error) {
	var ratings model.Ratings

	if err := repo.DB.Where("food_recipe_id = ?", recipeID).Find(&ratings).Error; err != nil {
		return nil, err
	}

	return ratings, nil
}

func (repo Repository) Create(rating *model.Rating) error {
	if err := repo.DB.Create(rating).First(&rating).Error; err != nil {
		return err
	}

	return nil
}
