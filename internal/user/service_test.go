package user_test

import (
	"reflect"
	"strings"
	"testing"
	"wongnok/internal/model"
	"wongnok/internal/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestNewService(t *testing.T) {
	t.Run("ShouldFillProperties", func(t *testing.T) {
		service := user.NewService(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(service))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})
}

type ServiceUpsertWithClaimsTestSuite struct {
	suite.Suite

	// Dependencies
	service user.IService
	repo    *MockIRepository

	// Mock data
	respGetByID model.User
	errGetByID  error
	errUpsert   error
}

// This will run before each test
func (suite *ServiceUpsertWithClaimsTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &user.Service{
		Repository: suite.repo,
	}

	suite.respGetByID = model.User{}
	suite.errGetByID = nil
	suite.errUpsert = nil

	suite.repo.On("GetByID", mock.Anything).Return(func(string) (model.User, error) {
		return suite.respGetByID, suite.errGetByID
	})

	suite.repo.On("Upsert", mock.Anything).Return(func(*model.User) error {
		return suite.errUpsert
	})
}

func (suite *ServiceUpsertWithClaimsTestSuite) TestReturnUserUpdated() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}

	expectedUser := model.User{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}

	user, err := suite.service.UpsertWithClaims(claims)
	suite.NoError(err)

	suite.Equal(expectedUser, user)
}

func (suite *ServiceUpsertWithClaimsTestSuite) TestErrorWhenClaimsValidated() {
	claims := model.Claims{}

	user, err := suite.service.UpsertWithClaims(claims)
	suite.Error(err)
	suite.True(strings.HasPrefix(err.Error(), "claims invalid"))

	suite.Empty(user)
}

func (suite *ServiceUpsertWithClaimsTestSuite) TestErrorWhenGetByID() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}

	suite.errGetByID = assert.AnError

	user, err := suite.service.UpsertWithClaims(claims)
	suite.ErrorIs(err, assert.AnError)
	suite.True(strings.HasPrefix(err.Error(), "find user"))

	suite.Empty(user)
}

func (suite *ServiceUpsertWithClaimsTestSuite) TestErrorWhenUpsert() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}

	suite.errUpsert = assert.AnError

	user, err := suite.service.UpsertWithClaims(claims)
	suite.ErrorIs(err, assert.AnError)
	suite.True(strings.HasPrefix(err.Error(), "upsert user"))

	suite.Empty(user)
}

func TestServiceUpsertWithClaims(t *testing.T) {
	suite.Run(t, new(ServiceUpsertWithClaimsTestSuite))
}

type ServiceGetByIDTestSuite struct {
	suite.Suite

	// Dependencies
	service user.IService
	repo    *MockIRepository

	// Mock data
	respGetByID model.User
	errGetByID  error
}

func (suite *ServiceGetByIDTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &user.Service{
		Repository: suite.repo,
	}
	suite.respGetByID = model.User{}

	suite.errGetByID = nil

	suite.repo.On("GetByID", mock.Anything).Return(func(string) (model.User, error) {
		return suite.respGetByID, suite.errGetByID
	})
}

func (suite *ServiceGetByIDTestSuite) TestGetByID() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}

	expectedUser := model.User{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	suite.respGetByID = expectedUser

	user, err := suite.service.GetByID(claims)

	suite.NoError(err)
	suite.Equal(expectedUser, user)
}

func (suite *ServiceGetByIDTestSuite) TestErrorWhenGetByID() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	suite.errGetByID = assert.AnError

	user, err := suite.service.GetByID(claims)

	suite.Error(err)
	suite.Equal(model.User{}, user)
}

func (suite *ServiceGetByIDTestSuite) TestNotFoundWhenGetByID() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	suite.errGetByID = gorm.ErrRecordNotFound

	user, err := suite.service.GetByID(claims)

	suite.NoError(err)
	suite.Equal(model.User{}, user)
}

func TestServiceGetByID(t *testing.T) {
	suite.Run(t, new(ServiceGetByIDTestSuite))
}

type ServiceGetRecipesTestSuite struct {
	suite.Suite

	// Dependencies
	service user.IService
	repo    *MockIRepository

	// Mock data
	respGetByID    model.User
	errGetByID     error
	respGetRecipes model.FoodRecipes
	errGetRecipes  error
}

func (suite *ServiceGetRecipesTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &user.Service{
		Repository: suite.repo,
	}

	suite.respGetByID = model.User{}
	suite.respGetRecipes = model.FoodRecipes{}
	suite.errGetByID = nil
	suite.errGetRecipes = nil

	suite.repo.On("GetByID", mock.Anything).Return(func(string) (model.User, error) {
		return suite.respGetByID, suite.errGetByID
	})

	suite.repo.On("GetRecipes", mock.Anything).Return(func(string) (model.FoodRecipes, error) {
		return suite.respGetRecipes, suite.errGetRecipes
	})
}

func (suite *ServiceGetRecipesTestSuite) TestGetRecipesResponseFoodRecipes() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}

	expectedFoodRecipes := model.FoodRecipes{
		model.FoodRecipe{
			Model: gorm.Model{ID: 1},
			Name:  "Omelet",
			Ratings: model.Ratings{
				model.Rating{
					Score: 5,
				},
				model.Rating{
					Score: 3,
				},
			},
		},
	}
	suite.respGetRecipes = expectedFoodRecipes

	foodRecipes, err := suite.service.GetRecipes("1", claims)

	suite.NoError(err)
	suite.Equal(expectedFoodRecipes, foodRecipes)
	suite.Equal(float64(4), foodRecipes[0].AverageRating)
	suite.repo.AssertCalled(suite.T(), "GetByID", "ID")
	suite.repo.AssertCalled(suite.T(), "GetRecipes", "1")
}

func (suite *ServiceGetRecipesTestSuite) TestGetRecipesResponseErrorGetByID() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	suite.errGetByID = assert.AnError
	foodRecipes, err := suite.service.GetRecipes("1", claims)
	suite.ErrorIs(err, assert.AnError)
	suite.True(strings.HasPrefix(err.Error(), "find user"))

	suite.Empty(foodRecipes)
}
func (suite *ServiceGetRecipesTestSuite) TestGetRecipesResponseErrorGetRecipes() {
	claims := model.Claims{
		ID:        "ID",
		FirstName: "FirstName",
		LastName:  "LastName",
	}
	suite.errGetRecipes = assert.AnError
	foodRecipes, err := suite.service.GetRecipes("1", claims)
	suite.ErrorIs(err, assert.AnError)
	suite.True(strings.HasPrefix(err.Error(), "get recipes"))

	suite.Empty(foodRecipes)
}

func TestServiceGetRecipes(t *testing.T) {
	suite.Run(t, new(ServiceGetRecipesTestSuite))
}
