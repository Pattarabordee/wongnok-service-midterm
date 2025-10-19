package user_test

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"
	"time"
	"wongnok/internal/model"
	"wongnok/internal/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	driver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestNewRepository(t *testing.T) {

	t.Run("ShouldFillProperties", func(t *testing.T) {
		repo := user.NewRepository(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(repo))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type RepositoryTestSuite struct {
	suite.Suite
	ctx       context.Context
	container *postgres.PostgresContainer
	db        *gorm.DB
	repo      user.IRepository
}

// This will run once before all tests in the suite
func (suite *RepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	container, err := postgres.Run(
		suite.ctx,
		"postgres:17-alpine",
		postgres.WithInitScripts(filepath.Join("../..", "tests", "init-db.sql")),
		postgres.WithDatabase("wongnok-test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(
				(5 * time.Second),
			),
		),
	)
	suite.NoError(err)
	suite.container = container
}

// This will run once after all tests in the suite
func (suite *RepositoryTestSuite) TearDownSuite() {
	err := suite.container.Terminate(suite.ctx)
	suite.NoError(err)
}

// This will run before each test
func (suite *RepositoryTestSuite) SetupTest() {
	conn, err := suite.container.ConnectionString(suite.ctx, "sslmode=disable")
	suite.NoError(err)

	db, err := gorm.Open(driver.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	suite.NoError(err)

	suite.repo = &user.Repository{
		DB: db,
	}

	suite.db = db
}

// This will run after each test
func (suite *RepositoryTestSuite) TearDownTest() {
	sqldb, _ := suite.db.DB()
	sqldb.Close()
}

// Extend
type RepositorGetByIDTestSuite struct {
	RepositoryTestSuite
}

func (suite *RepositorGetByIDTestSuite) TestReturnUserWithIDMatched() {
	expectedUser := model.User{
		ID:        "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
		FirstName: "Demo",
		LastName:  "Tester",
	}

	user, err := suite.repo.GetByID("38fa4e9e-27de-42d5-a70f-9f01d41f32c2")
	suite.NoError(err)

	user.CreatedAt, user.UpdatedAt = time.Time{}, time.Time{}

	suite.Equal(expectedUser, user)
}

func (suite *RepositorGetByIDTestSuite) TestReturnErrorRecordNotFound() {
	user, err := suite.repo.GetByID("")
	suite.ErrorIs(err, gorm.ErrRecordNotFound)

	suite.Empty(user)
}

func TestRepositoryGetByID(t *testing.T) {
	suite.Run(t, new(RepositorGetByIDTestSuite))
}

// Extend
type RepositorUpsertTestSuite struct {
	RepositoryTestSuite
}

func (suite *RepositorUpsertTestSuite) TestUserMustBeCreated() {
	user := model.User{
		ID:        "e30e4978-6e27-406e-861d-935d7a923037",
		FirstName: "Fake",
		LastName:  "Mocker",
	}

	expectedUser := model.User{
		ID:        "e30e4978-6e27-406e-861d-935d7a923037",
		FirstName: "Fake",
		LastName:  "Mocker",
	}

	err := suite.repo.Upsert(&user)
	suite.NoError(err)

	user.CreatedAt, user.UpdatedAt = time.Time{}, time.Time{}

	suite.Equal(expectedUser, user)
}

func (suite *RepositorUpsertTestSuite) TestUserMustBeUpdated() {
	user := model.User{
		ID:        "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
		FirstName: "Demo",
		LastName:  "Faker",
	}

	expectedUser := model.User{
		ID:        "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
		FirstName: "Demo",
		LastName:  "Faker",
	}

	err := suite.repo.Upsert(&user)
	suite.NoError(err)

	user.CreatedAt, user.UpdatedAt = time.Time{}, time.Time{}

	suite.Equal(expectedUser, user)
}

func (suite *RepositorUpsertTestSuite) TestErrorWhenUpsert() {
	sqldb, _ := suite.db.DB()
	sqldb.Close()

	err := suite.repo.Upsert(&model.User{})
	suite.EqualError(err, "sql: database is closed")
}

func TestRepositoryUpsert(t *testing.T) {
	suite.Run(t, new(RepositorUpsertTestSuite))
}

// Extend
type RepositoryGetRecipesTestSuite struct {
	RepositoryTestSuite
}

func (suite *RepositoryGetRecipesTestSuite) TestGetRecipesReturnFoodRecipes() {
	mockTime := time.Time{}
	expectedFoodRecipes := model.FoodRecipes{
		{
			Model:             gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
			Name:              "Omlet",
			Description:       "Eggs fried?",
			Ingredient:        "Eggs",
			Instruction:       "Cooking",
			ImageURL:          nil,
			CookingDurationID: 1,
			CookingDuration: model.CookingDuration{
				Model: gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
				Name:  "5 - 10",
			},
			DifficultyID: 1,
			Difficulty: model.Difficulty{
				Model: gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
				Name:  "Easy",
			},
			Ratings: model.Ratings{
				model.Rating{
					Model:        gorm.Model{ID: 1, CreatedAt: mockTime, UpdatedAt: mockTime},
					Score:        5,
					FoodRecipeID: 1,
					UserID:       "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
				},
				model.Rating{
					Model:        gorm.Model{ID: 2, CreatedAt: mockTime, UpdatedAt: mockTime},
					Score:        3,
					FoodRecipeID: 1,
					UserID:       "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
				},
			},
			AverageRating: 0,
			UserID:        "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
			User: model.User{
				ID:        "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
				FirstName: "Demo",
				LastName:  "Tester",
			},
		},
	}

	foodRecipes, err := suite.repo.GetRecipes("38fa4e9e-27de-42d5-a70f-9f01d41f32c2")
	suite.NoError(err)

	for i := range foodRecipes {
		foodRecipes[i].CreatedAt = mockTime
		foodRecipes[i].UpdatedAt = mockTime
		foodRecipes[i].CookingDuration.CreatedAt = mockTime
		foodRecipes[i].CookingDuration.UpdatedAt = mockTime
		foodRecipes[i].Difficulty.CreatedAt = mockTime
		foodRecipes[i].Difficulty.UpdatedAt = mockTime
		foodRecipes[i].Ratings[0].CreatedAt = mockTime
		foodRecipes[i].Ratings[0].UpdatedAt = mockTime
		foodRecipes[i].Ratings[1].CreatedAt = mockTime
		foodRecipes[i].Ratings[1].UpdatedAt = mockTime
		foodRecipes[i].User.CreatedAt = mockTime
		foodRecipes[i].User.UpdatedAt = mockTime
	}

	suite.Equal(expectedFoodRecipes, foodRecipes)
}

func TestRepositoryGetRecipes(t *testing.T) {
	suite.Run(t, new(RepositoryGetRecipesTestSuite))
}
