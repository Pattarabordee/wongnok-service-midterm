package foodrecipe_test

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"
	"time"
	"wongnok/internal/foodrecipe"
	"wongnok/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
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
		repo := foodrecipe.NewRepository(&gorm.DB{})

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
	repo      foodrecipe.IRepository
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

	suite.repo = &foodrecipe.Repository{
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
type RepositoryCreateTestSuite struct {
	RepositoryTestSuite
}

func (suite *RepositoryCreateTestSuite) TestReturnRecipePointerWhenCreated() {
	recipe := model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	err := suite.repo.Create(&recipe)
	suite.NoError(err)

	expectedRecipe := model.FoodRecipe{
		Model:             gorm.Model{ID: 2},
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		CookingDuration: model.CookingDuration{
			Model: gorm.Model{ID: 1},
			Name:  "5 - 10",
		},
		DifficultyID: 1,
		Difficulty: model.Difficulty{
			Model: gorm.Model{ID: 1},
			Name:  "Easy",
		},
		Ratings: model.Ratings{},
		UserID:  "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
		User: model.User{
			ID:        "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
			FirstName: "Demo",
			LastName:  "Tester",
		},
	}

	recipe.CookingDuration.CreatedAt, recipe.CookingDuration.UpdatedAt = time.Time{}, time.Time{}
	recipe.Difficulty.CreatedAt, recipe.Difficulty.UpdatedAt = time.Time{}, time.Time{}
	recipe.User.CreatedAt, recipe.User.UpdatedAt = time.Time{}, time.Time{}
	recipe.CreatedAt, recipe.UpdatedAt = time.Time{}, time.Time{}

	suite.Equal(expectedRecipe, recipe)
}

func (suite *RepositoryCreateTestSuite) TestErrorWhenDuplicateID() {
	recipe := model.FoodRecipe{
		Model:             gorm.Model{ID: 1},
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	err := suite.repo.Create(&recipe)
	suite.Error(err)
	suite.IsType(err, &pgconn.PgError{})

	suite.Equal("23505", err.(*pgconn.PgError).SQLState())
}

func TestRepositoryCreate(t *testing.T) {
	suite.Run(t, new(RepositoryCreateTestSuite))
}

// Extend
type RepositoryUpdateTestSuite struct {
	RepositoryTestSuite

	recipe model.FoodRecipe
}

func (suite *RepositoryUpdateTestSuite) SetupTest() {
	// Super
	suite.RepositoryTestSuite.SetupTest()

	suite.recipe = model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	err := suite.db.Create(&suite.recipe).Error
	suite.NoError(err)
}

func (suite *RepositoryUpdateTestSuite) TestUpdateRecipe() {
	recipe := model.FoodRecipe{
		Model: gorm.Model{ID: suite.recipe.ID},
		Name:  "Update Name",
	}

	err := suite.repo.Update(&recipe)
	suite.NoError(err)

	var result model.FoodRecipe
	err = suite.db.First(&result, suite.recipe.ID).Error
	suite.NoError(err)

	suite.Equal("Update Name", result.Name)
}

func (suite *RepositoryUpdateTestSuite) TestErrorWhenUpdate() {
	err := suite.repo.Update(&model.FoodRecipe{})
	suite.ErrorIs(err, gorm.ErrMissingWhereClause)
}

func TestRepositoryUpdate(t *testing.T) {
	suite.Run(t, new(RepositoryUpdateTestSuite))
}

// Extend
type RepositoryDeleteTestSuite struct {
	RepositoryTestSuite

	recipe model.FoodRecipe
}

func (suite *RepositoryDeleteTestSuite) SetupTest() {
	// Super
	suite.RepositoryTestSuite.SetupTest()

	suite.recipe = model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	err := suite.db.Create(&suite.recipe).Error
	suite.NoError(err)
}

func (suite *RepositoryDeleteTestSuite) TestDeleteRecipe() {
	err := suite.repo.Delete(int(suite.recipe.ID))
	suite.NoError(err)

	var result model.FoodRecipe
	err = suite.db.First(&result, suite.recipe.ID).Error
	suite.ErrorIs(err, gorm.ErrRecordNotFound)
}

func TestRepositoryDelete(t *testing.T) {
	suite.Run(t, new(RepositoryDeleteTestSuite))
}

type RepositoryGetTestSuite struct {
	RepositoryTestSuite

	recipe model.FoodRecipe
}

func (suite *RepositoryGetTestSuite) SetupTest() {
	// Super
	suite.RepositoryTestSuite.SetupTest()

	suite.recipe = model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	err := suite.db.Create(&suite.recipe).Error
	suite.NoError(err)
}

func (suite *RepositoryGetTestSuite) TearDownTest() {
	sqldb, _ := suite.db.DB()
	sqldb.Close()
}

func (suite *RepositoryGetTestSuite) TestGetRecipe() {

	foodRecipeQuery := model.FoodRecipeQuery{
		Search: "",
		Page:   1,
		Limit:  10,
	}

	response, err := suite.repo.Get(foodRecipeQuery)

	suite.NoError(err)
	suite.Equal(2, len(response))
	suite.Equal(uint(2), response[0].ID)
	suite.Equal(uint(1), response[1].ID)
}

func (suite *RepositoryGetTestSuite) TestGetRecipeLimit() {

	foodRecipeQuery := model.FoodRecipeQuery{
		Search: "",
		Page:   1,
		Limit:  1,
	}

	response, err := suite.repo.Get(foodRecipeQuery)

	suite.NoError(err)
	suite.Equal(1, len(response))
}

func (suite *RepositoryGetTestSuite) TestGetRecipeSearch() {

	foodRecipeQuery := model.FoodRecipeQuery{
		Search: "mle",
		Page:   1,
		Limit:  10,
	}

	response, err := suite.repo.Get(foodRecipeQuery)

	suite.NoError(err)
	suite.Equal(1, len(response))
	suite.Contains(response[0].Name, foodRecipeQuery.Search)
}

func TestRepositoryGet(t *testing.T) {
	suite.Run(t, new(RepositoryGetTestSuite))
}

type RepositoryCountTestSuite struct {
	RepositoryTestSuite

	recipe model.FoodRecipe
}

func (suite *RepositoryCountTestSuite) SetupTest() {
	// Super
	suite.RepositoryTestSuite.SetupTest()

	suite.recipe = model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	err := suite.db.Create(&suite.recipe).Error
	suite.NoError(err)
}

func (suite *RepositoryCountTestSuite) TestCount() {

	count, err := suite.repo.Count()
	suite.NoError(err)

	suite.Equal(int64(2), count)
}

func TestRepositoryCount(t *testing.T) {
	suite.Run(t, new(RepositoryCountTestSuite))
}

type RepositoryGetByIDTestSuite struct {
	RepositoryTestSuite

	recipe model.FoodRecipe
}

func (suite *RepositoryGetByIDTestSuite) SetupTest() {
	// Super
	suite.RepositoryTestSuite.SetupTest()

	suite.recipe = model.FoodRecipe{
		Name:              "Name",
		Description:       "Description",
		Ingredient:        "Ingredient",
		Instruction:       "Instruction",
		CookingDurationID: 1,
		DifficultyID:      1,
		UserID:            "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	err := suite.db.Create(&suite.recipe).Error
	suite.NoError(err)
}

func (suite *RepositoryGetByIDTestSuite) TestGetByID() {

	result, err := suite.repo.GetByID(1)
	suite.NoError(err)

	suite.Equal(uint(1), result.ID)
}

func (suite *RepositoryGetByIDTestSuite) TestErrorGetByID() {

	result, err := suite.repo.GetByID(99)

	suite.Equal(model.FoodRecipe{}, result)
	suite.ErrorIs(err, gorm.ErrRecordNotFound)
}

func TestRepositoryGetByID(t *testing.T) {
	suite.Run(t, new(RepositoryGetByIDTestSuite))
}
