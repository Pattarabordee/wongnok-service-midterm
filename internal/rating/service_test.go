package rating_test

import (
	"reflect"
	"testing"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"
	"wongnok/internal/rating"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestNewService(t *testing.T) {
	t.Run("Should Fill Properties", func(t *testing.T) {
		service := rating.NewService(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(service))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s should not be nil", field.Type().Name())
		}
	})
}

type ServiceGetRating struct {
	suite.Suite

	// Dependencies
	service rating.IService
	repo    *MockIRepository

	// Mock data
	respRepositoryGet model.Ratings
	errRepositoryGet  error
}

func (suite *ServiceGetRating) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &rating.Service{
		Repository: suite.repo,
	}

	suite.respRepositoryGet = model.Ratings{
		{
			Score:        4.5,
			FoodRecipeID: 1,

			UserID: "1a",
		},
		{
			Score:        5,
			FoodRecipeID: 1,

			UserID: "1a",
		},
		{
			Score:        3,
			FoodRecipeID: 1,

			UserID: "1a",
		},
	}

	suite.errRepositoryGet = nil

	// Mock the repository's Get method
	suite.repo.On("Get", mock.AnythingOfType("int")).Return(func(id int) (model.Ratings, error) {
		if id == 1 {
			return suite.respRepositoryGet, suite.errRepositoryGet
		}
		return nil, gorm.ErrRecordNotFound
	})
}

func (suite *ServiceGetRating) TestReturnRatingsWhenFound() {
	ratings, err := suite.service.Get(1)

	suite.repo.AssertCalled(suite.T(), "Get", 1)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), len(suite.respRepositoryGet), len(ratings))
	assert.Equal(suite.T(), suite.respRepositoryGet[0].Score, ratings[0].Score)
	assert.Equal(suite.T(), suite.respRepositoryGet[0].FoodRecipeID, ratings[0].FoodRecipeID)
	assert.Equal(suite.T(), suite.respRepositoryGet[1].Score, ratings[1].Score)
	assert.Equal(suite.T(), suite.respRepositoryGet[1].FoodRecipeID, ratings[1].FoodRecipeID)
	assert.Equal(suite.T(), suite.respRepositoryGet[2].Score, ratings[2].Score)
	assert.Equal(suite.T(), suite.respRepositoryGet[2].FoodRecipeID, ratings[2].FoodRecipeID)
}

func (suite *ServiceGetRating) TestReturnErrorWhenRepositoryNotFound() {
	ratings, err := suite.service.Get(2)

	// Expect call repo get with recipe ID 2
	suite.repo.AssertCalled(suite.T(), "Get", 2)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), ratings)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func TestServiceGetRating(t *testing.T) {
	suite.Run(t, new(ServiceGetRating))
}

type ServiceCreateRating struct {
	suite.Suite

	// Dependencies
	service rating.IService

	userService *MockIUserService

	repo *MockIRepository

	// Mock data
	errRepositoryCreate error

	errUserServiceGetByID error
	user                  model.User
}

func (suite *ServiceCreateRating) SetupTest() {
	suite.repo = new(MockIRepository)

	suite.userService = new(MockIUserService)

	suite.service = &rating.Service{
		Repository: suite.repo,

		UserService: suite.userService,
	}

	suite.user = model.User{
		ID:        "123abc",
		FirstName: "mock_FirstName_1",
		LastName:  "mock_LastName_1",
	}

	suite.errRepositoryCreate = nil

	suite.errUserServiceGetByID = nil

	suite.repo.On("Create", mock.AnythingOfType("*model.Rating")).Run(func(args mock.Arguments) {
		rating := args.Get(0).(*model.Rating)
		*rating = model.Rating{
			Score:        1,
			FoodRecipeID: 1,

			UserID: "1",
		}
	}).Return(func(*model.Rating) error {
		return suite.errRepositoryCreate
	})

	suite.userService.On("GetByID", mock.AnythingOfType("model.Claims")).Return(func(model.Claims) (model.User, error) {
		if suite.errUserServiceGetByID != nil {
			return model.User{}, suite.errUserServiceGetByID
		}
		return suite.user, suite.errUserServiceGetByID
	})

}

func (suite *ServiceCreateRating) TestReturnRatingWhenCreated() {
	request := dto.RatingRequest{
		Score: 1,
	}

	claims := model.Claims{
		ID:        "1",
		FirstName: "mock_FirstName",
		LastName:  "mock_LastName",
	}

	rating, err := suite.service.Create(request, 1, claims)

	suite.userService.AssertCalled(suite.T(), "GetByID", mock.AnythingOfType("model.Claims"))
	suite.repo.AssertCalled(suite.T(), "Create", mock.AnythingOfType("*model.Rating"))

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, int(rating.Score))
	assert.Equal(suite.T(), 1, int(rating.FoodRecipeID))
	assert.Equal(suite.T(), "1", rating.UserID)
}

func (suite *ServiceCreateRating) TestReturnErrorWhenRequestValidate() {
	request := dto.RatingRequest{
		Score: 0, // Invalid score
	}

	claims := model.Claims{
		ID:        "1234",
		FirstName: "mock_FirstName",
		LastName:  "mock_LastName",
	}

	rating, err := suite.service.Create(request, 1, claims)

	suite.userService.AssertNotCalled(suite.T(), "GetByID", mock.AnythingOfType("model.Claims"))

	suite.repo.AssertNotCalled(suite.T(), "Create", mock.AnythingOfType("*model.Rating"))

	assert.Equal(suite.T(), model.Rating{}, rating)
	assert.Equal(suite.T(), "request invalid: Key: 'RatingRequest.Score' Error:Field validation for 'Score' failed on the 'required' tag", err.Error())
}

func TestServiceCreateRating(t *testing.T) {
	suite.Run(t, new(ServiceCreateRating))
}
