package foodrecipe_test

import (
	"reflect"
	"strings"
	"testing"
	"wongnok/internal/foodrecipe"
	"wongnok/internal/global"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestNewService(t *testing.T) {

	t.Run("ShouldFillProperties", func(t *testing.T) {
		service := foodrecipe.NewService(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(service))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type ServiceCreateTestSuite struct {
	suite.Suite

	// Dependencies
	service foodrecipe.IService
	repo    *MockIRepository

	// Mock data
	errRepositoryCreate error
}

// This will run before each test
func (suite *ServiceCreateTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &foodrecipe.Service{
		Repository: suite.repo,
	}

	suite.errRepositoryCreate = nil

	suite.repo.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		recipe := args.Get(0).(*model.FoodRecipe)
		*recipe = model.FoodRecipe{
			Name:              "Name",
			Description:       "Description",
			Ingredient:        "Ingredient",
			Instruction:       "Instruction",
			CookingDurationID: 1,
			DifficultyID:      1,
			UserID:            "UID",
		}
	}).Return(func(*model.FoodRecipe) error {
		return suite.errRepositoryCreate
	})
}

func (suite *ServiceCreateTestSuite) TestReturnRecipeCreated() {
	claims := model.Claims{
		ID: "UID",
	}
	expectedRecipe := model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "UID",
	}

	recipe, err := suite.service.Create(
		dto.FoodRecipeRequest{
			Name:              "Name",
			Description:       "Description",
			Ingredient:        "Ingredient",
			Instruction:       "Instruction",
			CookingDurationID: 1,
			DifficultyID:      1,
		},
		claims,
	)
	suite.NoError(err)

	suite.Equal(expectedRecipe, recipe)
	suite.repo.AssertCalled(suite.T(), "Create", &model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "UID",
	})
}

func (suite *ServiceCreateTestSuite) TestErrorWhenRequestValidate() {
	recipe, err := suite.service.Create(dto.FoodRecipeRequest{}, model.Claims{})
	suite.Error(err)
	suite.True(strings.HasPrefix(err.Error(), "request invalid"))

	suite.Empty(recipe)
	suite.repo.AssertNotCalled(suite.T(), "Create")
}

func (suite *ServiceCreateTestSuite) TestErrorWhenRepositoryCreate() {
	suite.errRepositoryCreate = assert.AnError

	recipe, err := suite.service.Create(
		dto.FoodRecipeRequest{
			Name:              "Name",
			Description:       "Description",
			Ingredient:        "Ingredient",
			Instruction:       "Instruction",
			CookingDurationID: 1,
			DifficultyID:      1,
		},
		model.Claims{},
	)
	suite.ErrorIs(err, assert.AnError)

	suite.Empty(recipe)
}

func TestServiceCreateRecipe(t *testing.T) {
	suite.Run(t, new(ServiceCreateTestSuite))
}

type ServiceGetTestSuite struct {
	suite.Suite

	// Dependencies
	service foodrecipe.IService
	repo    *MockIRepository

	// Mock data
	respRepositoryCount int64
	errRepositoryCount  error
	respRepositoryGet   model.FoodRecipes
	errRepositoryGet    error
}

// This will run before each test
func (suite *ServiceGetTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &foodrecipe.Service{
		Repository: suite.repo,
	}

	suite.respRepositoryCount = 10
	suite.errRepositoryCount = nil
	suite.respRepositoryGet = model.FoodRecipes{
		{
			Name: "Name",
		},
	}
	suite.errRepositoryGet = nil

	suite.repo.On("Count").Return(func() (int64, error) {
		return suite.respRepositoryCount, suite.errRepositoryCount
	})

	suite.repo.On("Get", mock.AnythingOfType("model.FoodRecipeQuery")).Return(func(model.FoodRecipeQuery) (model.FoodRecipes, error) {
		return suite.respRepositoryGet, suite.errRepositoryGet
	})
}

func (suite *ServiceGetTestSuite) TestReturnRecipes() {

	foodRecipeQuery := model.FoodRecipeQuery{
		Search: "Name",
		Page:   1,
		Limit:  10,
	}
	recipes, total, err := suite.service.Get(foodRecipeQuery)

	suite.NoError(err)
	suite.Equal(model.FoodRecipes{
		{
			Name: "Name",
		},
	}, recipes)
	suite.Equal(int64(10), total)
}

func (suite *ServiceGetTestSuite) TestErrorWhenGet() {
	suite.respRepositoryGet = nil
	suite.errRepositoryGet = assert.AnError

	foodRecipeQuery := model.FoodRecipeQuery{
		Search: "Name",
		Page:   1,
		Limit:  10,
	}
	recipes, total, err := suite.service.Get(foodRecipeQuery)

	suite.ErrorIs(err, assert.AnError)

	suite.Empty(recipes)
	suite.Empty(total)
}

func (suite *ServiceGetTestSuite) TestErrorWhenCount() {
	suite.respRepositoryGet = nil
	suite.respRepositoryCount = 0
	suite.errRepositoryCount = assert.AnError

	foodRecipeQuery := model.FoodRecipeQuery{
		Search: "Name",
		Page:   1,
		Limit:  10,
	}

	recipes, total, err := suite.service.Get(foodRecipeQuery)

	suite.ErrorIs(err, assert.AnError)

	suite.Empty(recipes)
	suite.Empty(total)

	suite.repo.AssertNotCalled(suite.T(), "Get")
}

func TestServiceGetRecipes(t *testing.T) {
	suite.Run(t, new(ServiceGetTestSuite))
}

type ServiceGetByIDTestSuite struct {
	suite.Suite

	// Dependencies
	service foodrecipe.IService
	repo    *MockIRepository

	// Mock data
	respRepositoryGetByID model.FoodRecipe
	errRepositoryGetByID  error
}

func (suite *ServiceGetByIDTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &foodrecipe.Service{
		Repository: suite.repo,
	}

	suite.respRepositoryGetByID = model.FoodRecipe{
		Name: "Name",
	}
	suite.errRepositoryGetByID = nil

	suite.repo.On("GetByID", mock.AnythingOfType("int")).Return(func(id int) (model.FoodRecipe, error) {
		if id == 1 {
			return suite.respRepositoryGetByID, suite.errRepositoryGetByID
		}
		return model.FoodRecipe{}, gorm.ErrRecordNotFound
	})
}

func (suite *ServiceGetByIDTestSuite) TestReturnRecipeWhenFound() {
	recipe, err := suite.service.GetByID(1)
	suite.NoError(err)

	expectedRecipe := model.FoodRecipe{
		Name: "Name",
	}

	suite.Equal(expectedRecipe, recipe)
	suite.repo.AssertCalled(suite.T(), "GetByID", 1)
}

func (suite *ServiceGetByIDTestSuite) TestErrorWhenNotFound() {
	recipe, err := suite.service.GetByID(2)

	suite.ErrorIs(err, gorm.ErrRecordNotFound)
	suite.Empty(recipe)
	suite.repo.AssertCalled(suite.T(), "GetByID", 2)
}

func TestServiceGetRecipeByID(t *testing.T) {
	suite.Run(t, new(ServiceGetByIDTestSuite))
}

type ServiceUpdateTestSuite struct {
	suite.Suite

	// Dependencies
	service foodrecipe.IService
	repo    *MockIRepository

	// Mock data
	respGetByID          model.FoodRecipe
	errGetByID           error
	respRepositoryUpdate model.FoodRecipe
	errRepositoryUpdate  error
}

func (suite *ServiceUpdateTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &foodrecipe.Service{
		Repository: suite.repo,
	}

	suite.respGetByID = model.FoodRecipe{
		Model:             gorm.Model{ID: 1},
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "UID",
	}
	suite.errGetByID = nil
	suite.respRepositoryUpdate = model.FoodRecipe{
		Model:             gorm.Model{ID: 1},
		Name:              "NameUpdated",
		Description:       "DescriptionUpdated",
		Ingredient:        "IngredientUpdated",
		Instruction:       "InstructionUpdated",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "UID",
	}
	suite.errRepositoryUpdate = nil

	suite.repo.On("GetByID", mock.Anything).Return(func(int) (model.FoodRecipe, error) {
		return suite.respGetByID, suite.errGetByID
	})
	suite.repo.On("Update", mock.Anything).Run(func(args mock.Arguments) {
		recipe := args.Get(0).(*model.FoodRecipe)
		*recipe = suite.respRepositoryUpdate
	}).Return(func(*model.FoodRecipe) error {
		return suite.errRepositoryUpdate
	})
}

func (suite *ServiceUpdateTestSuite) TestReturnRecipeUpdated() {
	claims := model.Claims{
		ID: "UID",
	}
	expectedRecipe := model.FoodRecipe{
		Model:             gorm.Model{ID: 1},
		Name:              "NameUpdated",
		Description:       "DescriptionUpdated",
		Ingredient:        "IngredientUpdated",
		Instruction:       "InstructionUpdated",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "UID",
	}

	recipe, err := suite.service.Update(
		dto.FoodRecipeRequest{
			Name:              "NameUpdated",
			Description:       "DescriptionUpdated",
			Ingredient:        "IngredientUpdated",
			Instruction:       "InstructionUpdated",
			CookingDurationID: 1,
			DifficultyID:      1,
		},
		1,
		claims,
	)
	suite.NoError(err)

	suite.Equal(expectedRecipe, recipe)
	suite.repo.AssertCalled(suite.T(), "GetByID", 1)
	suite.repo.AssertCalled(suite.T(), "Update", &model.FoodRecipe{
		Model:             gorm.Model{ID: 1},
		Name:              "NameUpdated",
		Description:       "DescriptionUpdated",
		Ingredient:        "IngredientUpdated",
		Instruction:       "InstructionUpdated",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "UID",
	})
}

func (suite *ServiceUpdateTestSuite) TestErrorWhenRequestValidate() {
	recipe, err := suite.service.Update(dto.FoodRecipeRequest{}, 1, model.Claims{})
	suite.Error(err)
	suite.True(strings.HasPrefix(err.Error(), "request invalid"))

	suite.Empty(recipe)
	suite.repo.AssertNotCalled(suite.T(), "GetByID", mock.Anything)
	suite.repo.AssertNotCalled(suite.T(), "Update", mock.Anything)
}

func (suite *ServiceUpdateTestSuite) TestErrorWhenGetByID() {
	claims := model.Claims{
		ID: "UID",
	}

	suite.errGetByID = assert.AnError

	recipe, err := suite.service.Update(
		dto.FoodRecipeRequest{
			Name:              "NameUpdated",
			Description:       "DescriptionUpdated",
			Ingredient:        "IngredientUpdated",
			Instruction:       "InstructionUpdated",
			CookingDurationID: 1,
			DifficultyID:      1,
		},
		1,
		claims,
	)
	suite.ErrorIs(err, assert.AnError)
	suite.True(strings.HasPrefix(err.Error(), "find recipe"))

	suite.Empty(recipe)

	suite.repo.AssertNotCalled(suite.T(), "Update", mock.Anything)
}

func (suite *ServiceUpdateTestSuite) TestErrorForbidden() {
	claims := model.Claims{
		ID: "FAKE",
	}

	recipe, err := suite.service.Update(
		dto.FoodRecipeRequest{
			Name:              "NameUpdated",
			Description:       "DescriptionUpdated",
			Ingredient:        "IngredientUpdated",
			Instruction:       "InstructionUpdated",
			CookingDurationID: 1,
			DifficultyID:      1,
		},
		1,
		claims,
	)
	suite.ErrorIs(err, global.ErrForbidden)
	suite.Empty(recipe)

	suite.repo.AssertNotCalled(suite.T(), "Update", mock.Anything)
}

func (suite *ServiceUpdateTestSuite) TestErrorWhenRepositoryUpdate() {
	claims := model.Claims{
		ID: "UID",
	}

	suite.errRepositoryUpdate = assert.AnError

	recipe, err := suite.service.Update(
		dto.FoodRecipeRequest{
			Name:              "NameUpdated",
			Description:       "DescriptionUpdated",
			Ingredient:        "IngredientUpdated",
			Instruction:       "InstructionUpdated",
			CookingDurationID: 1,
			DifficultyID:      1,
		},
		1,
		claims,
	)
	suite.ErrorIs(err, assert.AnError)
	suite.True(strings.HasPrefix(err.Error(), "update recipe"))

	suite.Empty(recipe)
}

func TestServiceUpdateRecipe(t *testing.T) {
	suite.Run(t, new(ServiceUpdateTestSuite))
}

type ServiceDeleteTestSuite struct {
	suite.Suite

	// Dependencies
	service foodrecipe.IService
	repo    *MockIRepository

	// Mock data
	respGetByID         model.FoodRecipe
	errGetByID          error
	errRepositoryDelete error
}

func (suite *ServiceDeleteTestSuite) SetupTest() {
	suite.repo = new(MockIRepository)
	suite.service = &foodrecipe.Service{
		Repository: suite.repo,
	}

	suite.respGetByID = model.FoodRecipe{
		UserID: "UID",
	}
	suite.errGetByID = nil
	suite.errRepositoryDelete = nil

	suite.repo.On("GetByID", mock.Anything).Return(func(int) (model.FoodRecipe, error) {
		return suite.respGetByID, suite.errGetByID
	})
	suite.repo.On("Delete", mock.AnythingOfType("int")).Return(func(int) error {
		return suite.errRepositoryDelete
	})
}

func (suite *ServiceDeleteTestSuite) TestRecipeDeleted() {
	claims := model.Claims{
		ID: "UID",
	}

	err := suite.service.Delete(1, claims)
	suite.NoError(err)

	suite.repo.AssertCalled(suite.T(), "GetByID", 1)
	suite.repo.AssertCalled(suite.T(), "Delete", 1)
}

func (suite *ServiceDeleteTestSuite) TestErrorWhenGetByID() {
	claims := model.Claims{
		ID: "UID",
	}

	suite.errGetByID = assert.AnError

	err := suite.service.Delete(1, claims)
	suite.ErrorIs(err, assert.AnError)
	suite.True(strings.HasPrefix(err.Error(), "find recipe"))

	suite.repo.AssertNotCalled(suite.T(), "Delete", mock.Anything)
}

func (suite *ServiceDeleteTestSuite) TestErrorFoebidden() {
	claims := model.Claims{
		ID: "FAKE",
	}

	err := suite.service.Delete(1, claims)
	suite.ErrorIs(err, global.ErrForbidden)

	suite.repo.AssertNotCalled(suite.T(), "Delete", mock.Anything)
}

func (suite *ServiceDeleteTestSuite) TestErrorWhenRepositoryDelete() {
	claims := model.Claims{
		ID: "UID",
	}

	suite.errRepositoryDelete = assert.AnError

	err := suite.service.Delete(1, claims)
	suite.ErrorIs(err, assert.AnError)
}

func TestServiceDeleteRecipe(t *testing.T) {
	suite.Run(t, new(ServiceDeleteTestSuite))
}
