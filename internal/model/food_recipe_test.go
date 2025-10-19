package model_test

import (
	"testing"
	"time"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestFoodRecipeFromRequest(t *testing.T) {
	t.Run("ShouldSetFoodRecipeModel", func(t *testing.T) {
		imageURL := "ImageURL"

		request := dto.FoodRecipeRequest{
			Name:              "Name",
			Description:       "Description",
			Ingredient:        "Ingredient",
			Instruction:       "Instruction",
			ImageURL:          &imageURL,
			CookingDurationID: 1,
			DifficultyID:      1,
		}

		claims := model.Claims{
			ID: "UID",
		}

		var recipe model.FoodRecipe

		recipe = recipe.FromRequest(request, claims)

		expectedRecipe := model.FoodRecipe{
			Name:              "Name",
			Description:       "Description",
			Ingredient:        "Ingredient",
			Instruction:       "Instruction",
			ImageURL:          &imageURL,
			CookingDurationID: 1,
			DifficultyID:      1,
			UserID:            "UID",
		}

		assert.Equal(t, expectedRecipe, recipe)
	})

	t.Run("ShouldSetFoodRecipeModelWhenGormModelExists", func(t *testing.T) {
		imageURL := "ImageURL"

		request := dto.FoodRecipeRequest{
			Name:              "Name",
			Description:       "Description",
			Ingredient:        "Ingredient",
			Instruction:       "Instruction",
			ImageURL:          &imageURL,
			CookingDurationID: 1,
			DifficultyID:      1,
		}

		claims := model.Claims{
			ID: "UID",
		}

		recipe := model.FoodRecipe{
			Model: gorm.Model{ID: 1},
		}

		recipe = recipe.FromRequest(request, claims)

		expectedRecipe := model.FoodRecipe{
			Model:             gorm.Model{ID: 1},
			Name:              "Name",
			Description:       "Description",
			Ingredient:        "Ingredient",
			Instruction:       "Instruction",
			ImageURL:          &imageURL,
			CookingDurationID: 1,
			DifficultyID:      1,
			UserID:            "UID",
		}

		assert.Equal(t, expectedRecipe, recipe)
	})
}

func TestFoodRecipeToResponse(t *testing.T) {
	mockTime := time.Date(2025, 7, 9, 10, 10, 10, 0, time.Local)

	t.Run("ShouldReturnFoodRecipeResponse", func(t *testing.T) {
		imageURL := "ImageURL"

		recipe := model.FoodRecipe{
			Model:       gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
			Name:        "Name",
			Description: "Description",
			Ingredient:  "Ingredient",
			Instruction: "Instruction",
			ImageURL:    &imageURL,
			CookingDuration: model.CookingDuration{
				Model: gorm.Model{ID: 1},
				Name:  "CookingDuration",
			},
			Difficulty: model.Difficulty{
				Model: gorm.Model{ID: 2},
				Name:  "DifficultyName",
			},
			User: model.User{
				ID: "UID",
			},
		}

		expectedResponse := dto.FoodRecipeResponse{
			ID:          1,
			Name:        "Name",
			Description: "Description",
			Ingredient:  "Ingredient",
			Instruction: "Instruction",
			ImageURL:    &imageURL,
			CookingDuration: dto.CookingDurationResponse{
				ID:   1,
				Name: "CookingDuration",
			},
			Difficulty: dto.DifficultyResponse{
				ID:   2,
				Name: "DifficultyName",
			},
			User: dto.UserResponse{
				ID: "UID",
			},
			CreatedAt: mockTime,
			UpdatedAt: mockTime,
		}

		assert.Equal(t, expectedResponse, recipe.ToResponse())
	})
}

func TestFoodRecipesResponse(t *testing.T) {
	mockTime := time.Date(2025, 7, 9, 10, 10, 10, 0, time.Local)

	t.Run("ShouldReturnFoodRecipesResponse", func(t *testing.T) {
		imageURL := "ImageURL"

		recipes := model.FoodRecipes{
			{
				Model:       gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
				Name:        "Name",
				Description: "Description",
				Ingredient:  "Ingredient",
				Instruction: "Instruction",
				ImageURL:    &imageURL,
				CookingDuration: model.CookingDuration{
					Model: gorm.Model{ID: 1},
					Name:  "CookingDuration",
				},
				Difficulty: model.Difficulty{
					Model: gorm.Model{ID: 2},
					Name:  "DifficultyName",
				},
				User: model.User{
					ID: "UID",
				},
			},
		}

		expectedResponse := dto.FoodRecipesResponse{
			Total: 10,
			Results: []dto.FoodRecipeResponse{
				{
					ID:          1,
					Name:        "Name",
					Description: "Description",
					Ingredient:  "Ingredient",
					Instruction: "Instruction",
					ImageURL:    &imageURL,
					CookingDuration: dto.CookingDurationResponse{
						ID:   1,
						Name: "CookingDuration",
					},
					Difficulty: dto.DifficultyResponse{
						ID:   2,
						Name: "DifficultyName",
					},
					User: dto.UserResponse{
						ID: "UID",
					},
					CreatedAt: mockTime,
					UpdatedAt: mockTime,
				},
			},
		}

		assert.Equal(t, expectedResponse, recipes.ToResponse(10))
	})
}

func TestCalculateAverageRating(t *testing.T) {
	mockTime := time.Date(2025, 7, 9, 10, 10, 10, 0, time.Local)
	imageURL := "ImageURL"
	t.Run("ShouldCalculateAverageRating", func(t *testing.T) {
		recipe := model.FoodRecipe{
			Model:       gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
			Name:        "Name",
			Description: "Description",
			Ingredient:  "Ingredient",
			Instruction: "Instruction",
			ImageURL:    &imageURL,
			CookingDuration: model.CookingDuration{
				Model: gorm.Model{ID: 1},
				Name:  "CookingDuration",
			},
			Difficulty: model.Difficulty{
				Model: gorm.Model{ID: 2},
				Name:  "DifficultyName",
			},
			Ratings: model.Ratings{
				model.Rating{
					Score: 5,
				},
				model.Rating{
					Score: 4,
				},
			},
			User: model.User{
				ID: "UID",
			},
		}

		AvgRatingResult := recipe.CalculateAverageRating().AverageRating

		assert.Equal(t, float64(4.5), AvgRatingResult)
	})

	t.Run("ShouldCalculateAverageRatingWhenNotHaveRating", func(t *testing.T) {
		recipe := model.FoodRecipe{
			Model:       gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
			Name:        "Name",
			Description: "Description",
			Ingredient:  "Ingredient",
			Instruction: "Instruction",
			ImageURL:    &imageURL,
			CookingDuration: model.CookingDuration{
				Model: gorm.Model{ID: 1},
				Name:  "CookingDuration",
			},
			Difficulty: model.Difficulty{
				Model: gorm.Model{ID: 2},
				Name:  "DifficultyName",
			},
			User: model.User{
				ID: "UID",
			},
		}

		AvgRatingResult := recipe.CalculateAverageRating().AverageRating

		assert.Equal(t, float64(0), AvgRatingResult)
	})
}

func TestCalculateAverageRatings(t *testing.T) {
	mockTime := time.Date(2025, 7, 9, 10, 10, 10, 0, time.Local)
	imageURL := "ImageURL"
	t.Run("ShouldCalculateAverageRatings", func(t *testing.T) {
		recipes := model.FoodRecipes{
			{
				Model:       gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
				Name:        "Name",
				Description: "Description",
				Ingredient:  "Ingredient",
				Instruction: "Instruction",
				ImageURL:    &imageURL,
				CookingDuration: model.CookingDuration{
					Model: gorm.Model{ID: 1},
					Name:  "CookingDuration",
				},
				Difficulty: model.Difficulty{
					Model: gorm.Model{ID: 2},
					Name:  "DifficultyName",
				},
				Ratings: model.Ratings{
					model.Rating{
						Score: 5,
					},
					model.Rating{
						Score: 4,
					},
				},
				User: model.User{
					ID: "UID",
				},
			},
			{
				Model:       gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
				Name:        "Name",
				Description: "Description",
				Ingredient:  "Ingredient",
				Instruction: "Instruction",
				ImageURL:    &imageURL,
				CookingDuration: model.CookingDuration{
					Model: gorm.Model{ID: 1},
					Name:  "CookingDuration",
				},
				Difficulty: model.Difficulty{
					Model: gorm.Model{ID: 2},
					Name:  "DifficultyName",
				},
				Ratings: model.Ratings{
					model.Rating{
						Score: 3,
					},
					model.Rating{
						Score: 1,
					},
				},
				User: model.User{
					ID: "UID",
				},
			},
			{
				Model:       gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
				Name:        "Name",
				Description: "Description",
				Ingredient:  "Ingredient",
				Instruction: "Instruction",
				ImageURL:    &imageURL,
				CookingDuration: model.CookingDuration{
					Model: gorm.Model{ID: 1},
					Name:  "CookingDuration",
				},
				Difficulty: model.Difficulty{
					Model: gorm.Model{ID: 2},
					Name:  "DifficultyName",
				},
				User: model.User{
					ID: "UID",
				},
			},
		}

		result := recipes.CalculateAverageRatings()

		assert.Equal(t, float64(4.5), result[0].AverageRating)
		assert.Equal(t, float64(2), result[1].AverageRating)
		assert.Equal(t, float64(0), result[2].AverageRating)
	})
}
