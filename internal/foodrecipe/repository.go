package foodrecipe

import (
	"wongnok/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IRepository interface {
	Create(recipe *model.FoodRecipe) error
	Get(foodRecipeQuery model.FoodRecipeQuery) (model.FoodRecipes, error)
	Count() (int64, error)
	GetByID(id int) (model.FoodRecipe, error)
	Update(recipe *model.FoodRecipe) error
	Delete(id int) error
}

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) IRepository {
	return &Repository{
		DB: db,
	}
}

func (repo Repository) Create(recipe *model.FoodRecipe) error {
	return repo.DB.Preload(clause.Associations).Create(recipe).First(&recipe).Error
}

func (repo Repository) Get(query model.FoodRecipeQuery) (model.FoodRecipes, error) {
	var recipes = make(model.FoodRecipes, 0)

	offset := (query.Page - 1) * query.Limit
	db := repo.DB.Preload(clause.Associations)

	if query.Search != "" {
		db = db.Where("name LIKE ?", "%"+query.Search+"%").Or("description LIKE ?", "%"+query.Search+"%")
	}

	if err := db.Order("name asc").Limit(query.Limit).Offset(offset).Find(&recipes).Error; err != nil {
		return nil, err
	}

	return recipes, nil
}

func (repo Repository) Count() (int64, error) {
	var count int64

	if err := repo.DB.Model(&model.FoodRecipes{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (repo Repository) GetByID(id int) (model.FoodRecipe, error) {
	var recipe model.FoodRecipe

	if err := repo.DB.Preload(clause.Associations).First(&recipe, id).Error; err != nil {
		return model.FoodRecipe{}, err
	}

	return recipe, nil
}

func (repo Repository) Update(recipe *model.FoodRecipe) error {
	// update
	if err := repo.DB.Model(&recipe).Updates(recipe).Error; err != nil {
		return err
	}

	return repo.DB.Preload(clause.Associations).First(&recipe, recipe.ID).Error
}

func (repo Repository) Delete(id int) error {
	return repo.DB.Delete(&model.FoodRecipes{}, id).Error
}
