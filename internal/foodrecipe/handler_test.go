package foodrecipe_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"wongnok/internal/foodrecipe"
	"wongnok/internal/global"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestNewHandler(t *testing.T) {

	t.Run("ShouldFillProperties", func(t *testing.T) {
		handler := foodrecipe.NewHandler(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(handler))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type HandlerCreateTestSuite struct {
	suite.Suite

	// Dependencies
	handler foodrecipe.IHandler
	service *MockIService

	// Mock data
	claims            model.Claims
	respServiceCreate model.FoodRecipe
	errServiceCreate  error

	// Helper
	server func(payload io.Reader, claims *model.Claims) *httptest.ResponseRecorder
}

// This will run once before all tests in the suite
func (s *HandlerCreateTestSuite) SetupSuite() {
	// Gin testing mode
	gin.SetMode(gin.TestMode)
}

// This will run before each test
func (suite *HandlerCreateTestSuite) SetupTest() {
	// สร้าง mock ของ IService และ inject เข้า Handler
	suite.service = new(MockIService)
	suite.handler = foodrecipe.Handler{
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

		// Register route
		router.POST("/api/v1/food-recipes", suite.handler.Create)

		// ใช้ httptest.NewRecorder() จำลอง server และ request เพื่อตรวจสอบ HTTP response จริงๆ
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodPost,
			"/api/v1/food-recipes",
			payload,
		)
		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respServiceCreate = model.FoodRecipe{
		Name: "Name",
	}
	suite.errServiceCreate = nil

	suite.service.On("Create", mock.Anything, mock.Anything).Return(func(dto.FoodRecipeRequest, model.Claims) (model.FoodRecipe, error) {
		return suite.respServiceCreate, suite.errServiceCreate
	})
}

func (suite *HandlerCreateTestSuite) TestResponseRecipeWithStatusCode201() {
	payload := strings.NewReader(`{"name":"Name"}`)
	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	// Ensure close reader when terminated
	// Protect against memory leaks by closing the body after reading
	body := response.Result().Body
	defer body.Close()

	expectedResponse := dto.FoodRecipeResponse{
		Name: "Name",
	}
	expectedJson, _ := json.Marshal(expectedResponse)

	suite.Equal(http.StatusCreated, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())

	suite.service.AssertCalled(suite.T(), "Create", dto.FoodRecipeRequest{
		Name: "Name",
	}, model.Claims{
		ID: "UID",
	})
}

func (suite *HandlerCreateTestSuite) TestErrorWhenRequestInvalid() {
	payload := strings.NewReader(``)
	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":"EOF"}`, response.Body.String())
	suite.service.AssertNotCalled(suite.T(), "Create", mock.Anything)
}

func (suite *HandlerCreateTestSuite) TestErrorWhenServiceCreateRecipe() {
	suite.errServiceCreate = assert.AnError

	payload := strings.NewReader(`{"name":"Name"}`)
	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())
}

func (suite *HandlerCreateTestSuite) TestValidationErrorsErrorWhenServiceCreateRecipe() {
	suite.errServiceCreate = make(validator.ValidationErrors, 0)

	payload := strings.NewReader(`{"name":"Name"}`)
	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":""}`, response.Body.String())
}

func (suite *HandlerCreateTestSuite) TestErrorForbidden() {
	suite.errServiceCreate = global.ErrForbidden

	payload := strings.NewReader(`{"name":"Name"}`)
	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusForbidden, response.Code)
	suite.Equal(`{"message":"forbidden"}`, response.Body.String())
}

func (suite *HandlerCreateTestSuite) TestResponseStatusCode400() {
	response := suite.server(nil, nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Equal(`{"message":"Unauthorized"}`, response.Body.String())
	suite.service.AssertNotCalled(suite.T(), "GetRecipes")
}

func TestHandlerCreate(t *testing.T) {
	suite.Run(t, new(HandlerCreateTestSuite))
}

type HandlerGetTestSuite struct {
	suite.Suite

	// Dependencies
	handler foodrecipe.IHandler
	service *MockIService

	// Mock data
	respRecipesInServiceGet model.FoodRecipes
	respTotalInServiceGet   int64
	errServiceGet           error

	// Helper
	server func(payload io.Reader) *httptest.ResponseRecorder
}

// This will run once before all tests in the suite
func (sutie *HandlerGetTestSuite) SetupSuite() {
	// Gin testing mode
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerGetTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.handler = foodrecipe.Handler{
		Service: suite.service,
	}

	suite.server = func(payload io.Reader) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()
		router.GET("/api/v1/food-recipes", suite.handler.Get)

		// Recoder
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodGet,
			"/api/v1/food-recipes?search=name&page=1&limit=10",
			payload,
		)

		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respRecipesInServiceGet = model.FoodRecipes{
		{
			Name:        "Name",
			Description: "Description",
			Ingredient:  "Ingredient",
			Instruction: "Instruction",
			CookingDuration: model.CookingDuration{
				Model: gorm.Model{ID: 1},
				Name:  "CookingDurationName",
			},
			Difficulty: model.Difficulty{
				Model: gorm.Model{ID: 2},
				Name:  "DifficultyName",
			},
		},
	}
	suite.respTotalInServiceGet = 10
	suite.errServiceGet = nil

	suite.service.On("Get", mock.AnythingOfType("model.FoodRecipeQuery")).Return(func(model.FoodRecipeQuery) (model.FoodRecipes, int64, error) {
		return suite.respRecipesInServiceGet, suite.respTotalInServiceGet, suite.errServiceGet
	})
}

func (suite *HandlerGetTestSuite) TestResponseRecipesWithStatus200() {
	response := suite.server(nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	expectedResponse := dto.FoodRecipesResponse{
		Total: 10,
		Results: []dto.FoodRecipeResponse{
			{
				Name:        "Name",
				Description: "Description",
				Ingredient:  "Ingredient",
				Instruction: "Instruction",
				CookingDuration: dto.CookingDurationResponse{
					ID:   1,
					Name: "CookingDurationName",
				},
				Difficulty: dto.DifficultyResponse{
					ID:   2,
					Name: "DifficultyName",
				},
			},
		},
	}
	expectedJson, _ := json.Marshal(expectedResponse)

	suite.Equal(http.StatusOK, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())
}

func (suite *HandlerGetTestSuite) TestErrorWhenGetRecipes() {
	suite.errServiceGet = assert.AnError

	response := suite.server(nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())
}

func TestHandlerGet(t *testing.T) {
	suite.Run(t, new(HandlerGetTestSuite))
}

type HandlerGetByIDTestSuite struct {
	suite.Suite

	// Dependencies
	handler foodrecipe.IHandler
	service *MockIService

	// Mock data
	respRecipeInServiceGetByID model.FoodRecipe
	errServiceGetByID          error

	// Helper
	server func(payload io.Reader) *httptest.ResponseRecorder
}

func (suite *HandlerGetByIDTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerGetByIDTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.handler = foodrecipe.Handler{
		Service: suite.service,
	}

	suite.server = func(payload io.Reader) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()
		router.GET("/api/v1/food-recipes/:id", suite.handler.GetByID)

		// Recoder
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodGet,
			"/api/v1/food-recipes/1",
			payload,
		)
		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respRecipeInServiceGetByID = model.FoodRecipe{
		Model:       gorm.Model{ID: 1},
		Name:        "Name",
		Description: "Description",
		Ingredient:  "Ingredient",
		Instruction: "Instruction",
		CookingDuration: model.CookingDuration{
			Model: gorm.Model{ID: 1},
			Name:  "CookingDurationName",
		},
		Difficulty: model.Difficulty{
			Model: gorm.Model{ID: 1},
			Name:  "DifficultyName",
		},
	}

	suite.errServiceGetByID = nil
	suite.service.On("GetByID", mock.AnythingOfType("int")).Return(func(id int) (model.FoodRecipe, error) {
		if id == 1 {
			return suite.respRecipeInServiceGetByID, suite.errServiceGetByID
		}
		return model.FoodRecipe{}, gorm.ErrRecordNotFound
	})
}

func (suite *HandlerGetByIDTestSuite) TestResponseRecipeWithStatus200() {
	response := suite.server(nil)

	body := response.Result().Body
	defer body.Close()

	expectedResponse := dto.FoodRecipeResponse{
		ID:          1,
		Name:        "Name",
		Description: "Description",
		Ingredient:  "Ingredient",
		Instruction: "Instruction",
		CookingDuration: dto.CookingDurationResponse{
			ID:   1,
			Name: "CookingDurationName",
		},
		Difficulty: dto.DifficultyResponse{
			ID:   1,
			Name: "DifficultyName",
		},
	}
	expectedJson, _ := json.Marshal(expectedResponse)

	suite.Equal(http.StatusOK, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())
}

func (suite *HandlerGetByIDTestSuite) TestErrorWhenRecipeNotFound() {
	suite.errServiceGetByID = gorm.ErrRecordNotFound

	response := suite.server(nil)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"Recipe not found"}`, response.Body.String())
}

func (suite *HandlerGetByIDTestSuite) TestErrorWhenCallGetRecipe() {
	suite.errServiceGetByID = assert.AnError

	response := suite.server(nil)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())
}

func TestHandlerGetByID(t *testing.T) {
	suite.Run(t, new(HandlerGetByIDTestSuite))
}

type HandlerUpdateTestSuite struct {
	suite.Suite

	// Dependencies
	handler foodrecipe.IHandler
	service *MockIService

	// Mock data
	claims            model.Claims
	respServiceUpdate model.FoodRecipe
	errServiceUpdate  error

	// Helper
	server func(payload io.Reader, claims *model.Claims) *httptest.ResponseRecorder
}

func (suite *HandlerUpdateTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerUpdateTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.handler = foodrecipe.Handler{
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

		// Register route
		router.PUT("/api/v1/food-recipes/:id", suite.handler.Update)

		// Recoder
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodPut,
			"/api/v1/food-recipes/1",
			payload,
		)
		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respServiceUpdate = model.FoodRecipe{
		Name: "Name",
	}

	suite.errServiceUpdate = nil
	suite.service.On("Update", mock.Anything, mock.Anything, mock.Anything).Return(func(request dto.FoodRecipeRequest, id int, claims model.Claims) (model.FoodRecipe, error) {
		return suite.respServiceUpdate, suite.errServiceUpdate
	})
}

func (suite *HandlerUpdateTestSuite) TestResponseRecipeWithStatus200() {
	payload := strings.NewReader(`{"name":"Name"}`)
	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	body := response.Result().Body
	defer body.Close()

	expectedResponse := dto.FoodRecipeResponse{
		Name: "Name",
	}
	expectedJson, _ := json.Marshal(expectedResponse)

	suite.Equal(http.StatusOK, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())

	suite.service.AssertCalled(suite.T(), "Update", dto.FoodRecipeRequest{Name: "Name"}, 1, model.Claims{ID: "UID"})
}

func (suite *HandlerUpdateTestSuite) TestErrorWhenRequestInvalid() {
	payload := strings.NewReader(``)
	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":"EOF"}`, response.Body.String())

	suite.service.AssertNotCalled(suite.T(), "Update", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *HandlerUpdateTestSuite) TestErrorWhenCallServiceUpdateRecipe() {
	suite.errServiceUpdate = assert.AnError

	payload := strings.NewReader(`{}`)

	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())
}

func (suite *HandlerUpdateTestSuite) TestValidationErrorsWhenServiceUpdateRecipe() {
	suite.errServiceUpdate = make(validator.ValidationErrors, 0)

	payload := strings.NewReader(`{}`)

	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":""}`, response.Body.String())
}

func (suite *HandlerUpdateTestSuite) TestErrorFoebidden() {
	suite.errServiceUpdate = global.ErrForbidden

	payload := strings.NewReader(`{}`)

	claims := model.Claims{ID: "UID"}

	response := suite.server(payload, &claims)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusForbidden, response.Code)
	suite.Equal(`{"message":"forbidden"}`, response.Body.String())
}

func (suite *HandlerUpdateTestSuite) TestResponseStatusCode400() {
	response := suite.server(nil, nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Equal(`{"message":"Unauthorized"}`, response.Body.String())
	suite.service.AssertNotCalled(suite.T(), "GetRecipes")
}

func TestHandlerUpdate(t *testing.T) {
	suite.Run(t, new(HandlerUpdateTestSuite))
}

type HandlerDeleteTestSuite struct {
	suite.Suite

	// Dependencies
	handler foodrecipe.IHandler
	service *MockIService

	// Mock data
	claims           model.Claims
	errServiceDelete error

	// Helper
	server func(payload io.Reader, claims *model.Claims) *httptest.ResponseRecorder
}

func (suite *HandlerDeleteTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (suite *HandlerDeleteTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.handler = foodrecipe.Handler{
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

		// Register route
		router.DELETE("/api/v1/food-recipes/:id", suite.handler.Delete)

		// Recoder
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodDelete,
			"/api/v1/food-recipes/1",
			payload,
		)
		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.errServiceDelete = nil

	suite.service.On("Delete", mock.Anything, mock.Anything).Return(func(id int, claims model.Claims) error {
		return suite.errServiceDelete
	})
}

func (suite *HandlerDeleteTestSuite) TestResponseWithStatus200() {
	claims := model.Claims{ID: "UID"}
	response := suite.server(nil, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusOK, response.Code)
	suite.Equal(`{"message":"Recipe deleted successfully"}`, response.Body.String())

	suite.service.AssertCalled(suite.T(), "Delete", 1, model.Claims{ID: "UID"})
}

func (suite *HandlerDeleteTestSuite) TestErrorWhenCallServiceDeleteRecipe() {
	suite.errServiceDelete = assert.AnError

	claims := model.Claims{ID: "UID"}
	response := suite.server(nil, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())
}

func (suite *HandlerDeleteTestSuite) TestErrorFoebidden() {
	suite.errServiceDelete = global.ErrForbidden

	claims := model.Claims{ID: "UID"}
	response := suite.server(nil, &claims)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusForbidden, response.Code)
	suite.Equal(`{"message":"forbidden"}`, response.Body.String())
}

func (suite *HandlerDeleteTestSuite) TestResponseStatusCode400() {
	response := suite.server(nil, nil)

	// Ensure close reader when terminated
	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusUnauthorized, response.Code)
	suite.Equal(`{"message":"Unauthorized"}`, response.Body.String())
	suite.service.AssertNotCalled(suite.T(), "GetRecipes")
}

func TestHandlerDelete(t *testing.T) {
	suite.Run(t, new(HandlerDeleteTestSuite))
}
