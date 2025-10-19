package model_test

import (
	"testing"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/stretchr/testify/assert"
)

func TestRatingFromRequest(t *testing.T) {
	t.Run("ShouldSetFromRequestModel", func(t *testing.T) {
		request := dto.RatingRequest{
			Score: 4.5,
		}

		var rating model.Rating
		rating = rating.FromRequest(request)

		expectRating := model.Rating{
			Score: 4.5,
		}

		assert.Equal(t, expectRating.Score, rating.Score)
	})
}

func TestRatingToResponse(t *testing.T) {
	t.Run("ShouldSetToResponseModel", func(t *testing.T) {
		rating := model.Rating{
			Score:        4.5,
			FoodRecipeID: 1,
		}

		response := rating.ToResponse()

		expectRating := dto.RatingResponse{
			Score:        4.5,
			FoodRecipeID: 1,
		}

		assert.Equal(t, expectRating.Score, response.Score)
		assert.Equal(t, expectRating.FoodRecipeID, response.FoodRecipeID)
	})
}

func TestRatingsToResponse(t *testing.T) {
	t.Run("ShouldSetRatingsToResponseModel", func(t *testing.T) {
		ratings := model.Ratings{
			{Score: 4.5, FoodRecipeID: 1, UserID: "1a"},
			{Score: 3.0, FoodRecipeID: 2, UserID: "1a"},
		}

		response := ratings.ToResponse()

		expectRatings := dto.RatingsResponse{
			Results: []dto.RatingResponse{
				{Score: 4.5, FoodRecipeID: 1},
				{Score: 3.0, FoodRecipeID: 2},
			},
		}

		assert.Len(t, response.Results, 2)
		assert.Equal(t, expectRatings.Results[0].Score, response.Results[0].Score)
		assert.Equal(t, expectRatings.Results[1].Score, response.Results[1].Score)
	})
}
