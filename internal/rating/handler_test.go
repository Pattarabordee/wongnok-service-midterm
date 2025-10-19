package rating_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"
	"wongnok/internal/rating"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestNewHandler(t *testing.T) {

	t.Run("ShouldFillProperties", func(t *testing.T) {
		handler := rating.NewHandler(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(handler))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type HandlerCreateRatingTestSuite struct {
	suite.Suite

	// Dependencies
	handler rating.IHandler
	service *MockIService

	// Mock data

	claims model.Claims

	respServiceCreate model.Rating
	errServiceCreate  error

	// Helper
	server func(payload io.Reader, claims *model.Claims) *httptest.ResponseRecorder
}

// This will run once before all tests in the suite
func (suite *HandlerCreateRatingTestSuite) SetupSuite() {
	// Gin testing mode
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerCreateRatingTestSuite) SetupTest() {
	// สร้าง mock ของ IService และ inject เข้า Handler
	suite.service = new(MockIService)
	suite.handler = rating.Handler{
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

		router.POST("/api/v1/food-recipes/:id/ratings", suite.handler.Create)

		// ใช้ httptest.NewRecorder() จำลอง server และ request เพื่อตรวจสอบ HTTP response จริงๆ
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodPost,
			"/api/v1/food-recipes/1/ratings",
			payload,
		)

		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respServiceCreate = model.Rating{
		Score:        5,
		FoodRecipeID: 1,
	}

	suite.errServiceCreate = nil
	suite.service.On("Create",
		mock.AnythingOfType("dto.RatingRequest"),
		mock.AnythingOfType("int"),

		mock.AnythingOfType("model.Claims"),
	).Return(func(request dto.RatingRequest, id int, claim model.Claims) (model.Rating, error) {

		if id == 1 {
			return suite.respServiceCreate, suite.errServiceCreate
		}
		return model.Rating{}, assert.AnError
	})
}

func (suite *HandlerCreateRatingTestSuite) TestResponseRatingWithStatusCode201() {
	payload := strings.NewReader(`{"score": 5}`)

	claims := model.Claims{ID: "UID"}
	response := suite.server(payload, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	expectedBody := dto.RatingResponse{
		Score:        5,
		FoodRecipeID: 1,
	}

	expectedJson, _ := json.Marshal(expectedBody)

	suite.Equal(http.StatusCreated, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())
	suite.service.AssertCalled(suite.T(), "Create", dto.RatingRequest{
		Score: 5,
	}, 1, model.Claims{
		ID: "UID",
	})
}

func (suite *HandlerCreateRatingTestSuite) TestResponseErrorWhenRequestInvalid() {
	payload := strings.NewReader(``)

	claims := model.Claims{ID: "UID"}
	response := suite.server(payload, &claims)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":"EOF"}`, response.Body.String())
	suite.service.AssertNotCalled(suite.T(), "Create", mock.AnythingOfType("dto.RatingRequest"), mock.AnythingOfType("int"))
}

func (suite *HandlerCreateRatingTestSuite) TestValidationErrorsWhenServiceCreateRating() {
	suite.errServiceCreate = make(validator.ValidationErrors, 0)

	payload := strings.NewReader(`{"score": 5}`)

	claims := model.Claims{ID: "UID"}
	response := suite.server(payload, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":""}`, response.Body.String())
}

func (suite *HandlerCreateRatingTestSuite) TestResponseErrorStatusCode400() {
	payload := strings.NewReader(`{"score": 5}`)
	response := suite.server(payload, nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Equal(`{"message":"Unauthorized"}`, response.Body.String())
	suite.service.AssertNotCalled(suite.T(), "GetRecipes")
}

func TestHandlerCreateRating(t *testing.T) {
	suite.Run(t, new(HandlerCreateRatingTestSuite))
}

type HandlerGetRatingsTestSuite struct {
	suite.Suite

	// Dependencies
	handler rating.IHandler
	service *MockIService

	// Mock data
	respServiceGet model.Ratings
	errServiceGet  error

	// Helper
	server func(payload io.Reader) *httptest.ResponseRecorder
}

// This will run once before all tests in the suite
func (suite *HandlerGetRatingsTestSuite) SetupSuite() {
	// Gin testing mode
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerGetRatingsTestSuite) SetupTest() {
	// สร้าง mock ของ IService และ inject เข้า Handler
	suite.service = new(MockIService)
	suite.handler = rating.Handler{
		Service: suite.service,
	}

	suite.server = func(payload io.Reader) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()
		router.GET("/api/v1/food-recipes/:id/ratings", suite.handler.Get)

		// Recoder
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodGet,
			"/api/v1/food-recipes/1/ratings",
			payload,
		)
		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respServiceGet = model.Ratings{
		model.Rating{
			Score:        5,
			FoodRecipeID: 1,

			UserID: "1a",
		},
		model.Rating{
			Score:        4,
			FoodRecipeID: 1,

			UserID: "1a",
		},
	}

	suite.errServiceGet = nil
	suite.service.On("Get", mock.AnythingOfType("int")).Return(func(id int) (model.Ratings, error) {
		if id == 1 {
			return suite.respServiceGet, suite.errServiceGet
		}
		return model.Ratings{}, assert.AnError
	})
}

func (suite *HandlerGetRatingsTestSuite) TestResponseRatingsWithStatusCode200() {
	response := suite.server(nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	expectedResponse := dto.RatingsResponse{
		Results: []dto.RatingResponse{

			{Score: 5, FoodRecipeID: 1, UserID: "1a"},
			{Score: 4, FoodRecipeID: 1, UserID: "1a"},
		},
	}

	expectedJson, _ := json.Marshal(expectedResponse)

	suite.Equal(http.StatusOK, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())
	suite.service.AssertCalled(suite.T(), "Get", 1)
}

func (suite *HandlerGetRatingsTestSuite) TestResponseErrorWhenRecipeNotFound() {
	suite.errServiceGet = gorm.ErrRecordNotFound

	response := suite.server(nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"Rating not found"}`, response.Body.String())
}

func TestHandlerGetRatings(t *testing.T) {
	suite.Run(t, new(HandlerGetRatingsTestSuite))
}
