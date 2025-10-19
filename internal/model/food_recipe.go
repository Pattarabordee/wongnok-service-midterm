package model

import (
	"wongnok/internal/model/dto"

	"gorm.io/gorm"
)

type FoodRecipe struct {
	gorm.Model
	Name              string
	Description       string
	Ingredient        string
	Instruction       string
	ImageURL          *string
	CookingDurationID uint
	CookingDuration   CookingDuration
	DifficultyID      uint
	Difficulty        Difficulty
	Ratings           Ratings
	AverageRating     float64 `gorm:"-"`
	UserID            string
	User              User
}

func (recipe FoodRecipe) FromRequest(request dto.FoodRecipeRequest, claims Claims) FoodRecipe {
	return FoodRecipe{
		Model:             recipe.Model,
		Name:              request.Name,
		Description:       request.Description,
		Ingredient:        request.Ingredient,
		Instruction:       request.Instruction,
		ImageURL:          request.ImageURL,
		CookingDurationID: request.CookingDurationID,
		DifficultyID:      request.DifficultyID,
		UserID:            claims.ID,
	}
}

func (recipe FoodRecipe) ToResponse() dto.FoodRecipeResponse {
	return dto.FoodRecipeResponse{
		ID:          recipe.ID,
		Name:        recipe.Name,
		Description: recipe.Description,
		Ingredient:  recipe.Ingredient,
		Instruction: recipe.Instruction,
		ImageURL:    recipe.ImageURL,
		CookingDuration: dto.CookingDurationResponse{
			ID:   recipe.CookingDuration.ID,
			Name: recipe.CookingDuration.Name,
		},
		Difficulty: dto.DifficultyResponse{
			ID:   recipe.Difficulty.ID,
			Name: recipe.Difficulty.Name,
		},
		AverageRating: recipe.AverageRating,
		User:          recipe.User.ToResponse(),
		CreatedAt:     recipe.CreatedAt,
		UpdatedAt:     recipe.UpdatedAt,
	}
}

type FoodRecipes []FoodRecipe

func (recipes FoodRecipes) ToResponse(total int64) dto.FoodRecipesResponse {
	var results = make([]dto.FoodRecipeResponse, 0)

	for _, recipe := range recipes {
		results = append(results, recipe.ToResponse())
	}

	return dto.FoodRecipesResponse{
		Total:   total,
		Results: results,
	}
}

func (recipe FoodRecipe) CalculateAverageRating() FoodRecipe {
	if len(recipe.Ratings) > 0 {
		var totalRating float64
		for _, rating := range recipe.Ratings {
			totalRating += rating.Score
		}
		recipe.AverageRating = totalRating / float64(len(recipe.Ratings))
	} else {
		recipe.AverageRating = 0
	}
	return recipe
}

func (recipes FoodRecipes) CalculateAverageRatings() FoodRecipes {
	for i, recipe := range recipes {
		if len(recipe.Ratings) > 0 {
			var totalRating float64
			for _, rating := range recipe.Ratings {
				totalRating += rating.Score
			}
			recipes[i].AverageRating = totalRating / float64(len(recipe.Ratings))
		} else {
			recipes[i].AverageRating = 0
		}
	}
	return recipes
}

type FoodRecipeQuery struct {
	Search string `form:"search"`
	Page   int    `form:"page" binding:"required,min=1"`  // page number for pagination
	Limit  int    `form:"limit" binding:"required,min=1"` // number of items per page
}
