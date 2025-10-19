package dto

import "time"

type FoodRecipeRequest struct {
	Name              string  `validate:"required"`
	Description       string  `validate:"required"`
	Ingredient        string  `validate:"required"`
	Instruction       string  `validate:"required"`
	ImageURL          *string `validate:"omitempty,url"`
	CookingDurationID uint    `validate:"required"`
	DifficultyID      uint    `validate:"required"`
}

type FoodRecipeResponse struct {
	ID              uint                    `json:"id"`
	Name            string                  `json:"name"`
	Description     string                  `json:"description"`
	Ingredient      string                  `json:"ingredient"`
	Instruction     string                  `json:"instruction"`
	ImageURL        *string                 `json:"imageUrl,omitempty"`
	CookingDuration CookingDurationResponse `json:"cookingDuration"`
	Difficulty      DifficultyResponse      `json:"difficulty"`
	CreatedAt       time.Time               `json:"createdAt"`
	UpdatedAt       time.Time               `json:"updatedAt"`
	AverageRating   float64                 `json:"averageRating"`
	User            UserResponse            `json:"user"`
}

type FoodRecipesResponse BaseListResponse[[]FoodRecipeResponse]
