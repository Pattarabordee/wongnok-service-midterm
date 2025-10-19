package user_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"
	"wongnok/internal/user"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestNewHandler(t *testing.T) {

	t.Run("ShouldFillProperties", func(t *testing.T) {
		handler := user.NewHandler(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(handler))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type HandlerGetRecipesTestSuite struct {
	suite.Suite

	// Dependencies
	handler user.IHandler
	service *MockIService

	// Mock data
	mockTime              time.Time
	respServiceGetRecipes model.FoodRecipes
	errServiceGetRecipes  error

	// Helper
	server func(payload io.Reader, claims *model.Claims) *httptest.ResponseRecorder
}

// This will run once before all tests in the suite
func (suite *HandlerGetRecipesTestSuite) SetupSuite() {
	// Gin testing mode
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerGetRecipesTestSuite) SetupTest() {
	// สร้าง mock ของ IService และ inject เข้า Handler
	suite.service = new(MockIService)
	suite.handler = user.Handler{
		Service: suite.service,
	}

	suite.server = func(payload io.Reader, claims *model.Claims) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()

		// Set context
		router.Use(func(ctx *gin.Context) {
			if claims != nil {
				ctx.Set("claims", *claims)
			}
		})

		router.GET("/api/v1/users/:id/food-recipes", suite.handler.GetRecipes)

		// ใช้ httptest.NewRecorder() จำลอง server และ request เพื่อตรวจสอบ HTTP response จริงๆ
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodGet,
			"/api/v1/users/1/food-recipes",
			payload,
		)

		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respServiceGetRecipes = model.FoodRecipes{
		{
			Model:       gorm.Model{ID: 1, CreatedAt: suite.mockTime, UpdatedAt: suite.mockTime},
			Name:        "Name",
			Description: "Description",
			Ingredient:  "Ingredient",
			Instruction: "Instruction",
			ImageURL:    nil,
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

	suite.errServiceGetRecipes = nil
	suite.service.On("GetRecipes",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("model.Claims"),
	).Return(func(userID string, claim model.Claims) (model.FoodRecipes, error) {
		if userID == "1" {
			return suite.respServiceGetRecipes, suite.errServiceGetRecipes
		}
		return model.FoodRecipes{}, assert.AnError
	})
}

func (suite *HandlerGetRecipesTestSuite) TestResponseFoodRecipesWithStatusCode200() {
	claims := model.Claims{ID: "UID"}
	response := suite.server(nil, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	expectedResponse := dto.FoodRecipesResponse{
		Total: 1,
		Results: []dto.FoodRecipeResponse{
			{
				ID:          1,
				Name:        "Name",
				Description: "Description",
				Ingredient:  "Ingredient",
				Instruction: "Instruction",
				ImageURL:    nil,
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
				CreatedAt: suite.mockTime,
				UpdatedAt: suite.mockTime,
			},
		},
	}

	expectedJson, _ := json.Marshal(expectedResponse)

	suite.Equal(http.StatusOK, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())
	suite.service.AssertCalled(suite.T(), "GetRecipes", "1", model.Claims{ID: "UID"})
}

func (suite *HandlerGetRecipesTestSuite) TestResponseFoodRecipesWithStatusCode400() {
	response := suite.server(nil, nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Equal(`{"message":"Unauthorized"}`, response.Body.String())
	suite.service.AssertNotCalled(suite.T(), "GetRecipes")
}

func (suite *HandlerGetRecipesTestSuite) TestResponseErrorWhenGetRecipes() {
	suite.errServiceGetRecipes = assert.AnError

	claims := model.Claims{ID: "UID"}

	response := suite.server(nil, &claims)

	suite.Equal(http.StatusInternalServerError, response.Code)
}

func TestHandlerGetRecipes(t *testing.T) {
	suite.Run(t, new(HandlerGetRecipesTestSuite))
}
