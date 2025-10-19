package rating

import (
	"errors"
	"net/http"
	"strconv"
	"wongnok/internal/helper"
	"wongnok/internal/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type IHandler interface {
	Get(ctx *gin.Context)
	Create(ctx *gin.Context)
}

type Handler struct {
	Service IService
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		Service: NewService(db),
	}
}

// Get godoc
// @Summary Get ratings
// @Description Get ratings for a food recipe by ID
// @Tags ratings
// @Accept json
// @Produce json
// @Param id path string false "Food Recipe ID"
// @Success 200 {object} dto.FoodRecipesResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/food-recipes/{id}/ratings [get]
func (handler Handler) Get(ctx *gin.Context) {
	var id int

	pathParam := ctx.Param("id")
	if pathParam != "" {
		if parsed, err := strconv.Atoi(pathParam); err == nil && parsed > 0 {
			id = parsed
		}
	}

	ratings, err := handler.Service.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Rating not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ratings.ToResponse())
}

// Create godoc
// @Summary Create a rating
// @Description Create a new rating for a food recipe by ID
// @Tags ratings
// @Accept json
// @Produce json
// @Param id path string true "Food Recipe ID"
// @Param request body dto.RatingRequest true "Rating Request"
// @Success 201 {object} dto.RatingResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/v1/food-recipes/{id}/ratings [post]
func (handler Handler) Create(ctx *gin.Context) {
	var request dto.RatingRequest
	var id int

	pathParam := ctx.Param("id")
	if pathParam != "" {
		if parsed, err := strconv.Atoi(pathParam); err == nil && parsed > 0 {
			id = parsed
		}
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	claims, err := helper.DecodeClaims(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}

	rating, err := handler.Service.Create(request, id, claims)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.As(err, &validator.ValidationErrors{}) {
			statusCode = http.StatusBadRequest
		}

		ctx.JSON(statusCode, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, rating.ToResponse())
}
