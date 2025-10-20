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

	LikeRecipe(userID string, recipeID int) error
	UnlikeRecipe(userID string, recipeID int) error
	HasUserLoved(userID string, recipeID int) (bool, error)
	GetLovedRecipesByUser(userID string) ([]model.FoodRecipe, error)
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

func (repo Repository) LikeRecipe(userID string, recipeID int) error {
	love := model.RecipeLove{UserID: userID, RecipeID: recipeID}
	return repo.DB.Create(&love).Error
}

func (repo Repository) UnlikeRecipe(userID string, recipeID int) error {
	return repo.DB.Where("user_id = ? AND recipe_id = ?", userID, recipeID).Delete(&model.RecipeLove{}).Error
}

func (repo Repository) HasUserLoved(userID string, recipeID int) (bool, error) {
	var count int64
	err := repo.DB.Model(&model.RecipeLove{}).
		Where("user_id = ? AND recipe_id = ?", userID, recipeID).
		Count(&count).Error
	return count > 0, err
}

func (repo Repository) GetLovedRecipesByUser(userID string) ([]model.FoodRecipe, error) {
	var recipes []model.FoodRecipe
	err := repo.DB.
		Joins("JOIN recipe_loves ON recipe_loves.recipe_id = food_recipes.id").
		Joins("JOIN difficulties ON difficulties.id = food_recipes.difficulty_id").
		Joins("JOIN cooking_durations ON cooking_durations.id = food_recipes.cooking_duration_id").
		Joins("JOIN users ON users.id = food_recipes.user_id").
		Where("recipe_loves.user_id = ?", userID).
		Select("food_recipes.*, difficulties.id as difficulty__id, difficulties.name as difficulty__name, cooking_durations.id as cooking_duration__id, cooking_durations.name as cooking_duration__name, users.id as user__id, users.first_name as user__first_name, users.last_name as user__last_name").
		Find(&recipes).Error
	return recipes, err
}
