package rating_test

import (
	"context"
	"path/filepath"
	"reflect"
	"testing"
	"time"
	"wongnok/internal/model"
	"wongnok/internal/rating"

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
		repo := rating.NewRepository(&gorm.DB{})

		value := reflect.Indirect(reflect.ValueOf(repo))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})

}

type RepositoryTestSuite struct {
	suite.Suite
	ctx        context.Context
	container  *postgres.PostgresContainer
	db         *gorm.DB
	repository rating.IRepository
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

	suite.repository = &rating.Repository{
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
type RepositoryCreateRatingTestSuite struct {
	RepositoryTestSuite
}

func (suite *RepositoryCreateRatingTestSuite) TestReturnRatingPointerWhenRatingCreated() {
	rating := model.Rating{
		Score:        5,
		FoodRecipeID: 1,
		UserID:       "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	err := suite.repository.Create(&rating)
	suite.NoError(err)

	// GORM จะ auto-assign ID ที่ database generate ให้
	expectedRating := model.Rating{
		Model:        gorm.Model{ID: 3}, // ← Database auto-increment มี 1 อยู่แล้วตอน initial data
		Score:        5,
		FoodRecipeID: 1,
		UserID:       "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
	}

	// Ignore CreatedAt field for comparison
	rating.CreatedAt = time.Time{}
	rating.UpdatedAt = time.Time{}

	suite.Equal(expectedRating, rating)
}

func (suite *RepositoryCreateRatingTestSuite) TestErrorWhenDuplicateID() {
	rating := model.Rating{
		Model:        gorm.Model{ID: 1},
		Score:        5,
		FoodRecipeID: 1,
		UserID:       "1a",
	}

	err := suite.repository.Create(&rating)
	suite.Error(err)

	// SQL state 23505 is พยายามจะ insert ข้อมูลที่มี primary key ซ้ำ
	suite.Equal("23505", err.(*pgconn.PgError).SQLState())
}

func TestRepositoryCreateRating(t *testing.T) {
	suite.Run(t, new(RepositoryCreateRatingTestSuite))
}

type RepositoryGetRatingTestSuite struct {
	RepositoryTestSuite
}

func (suite *RepositoryGetRatingTestSuite) TestReturnRatings() {
	recipeID := 1

	result, err := suite.repository.Get(recipeID)
	suite.NoError(err)

	suite.NotEmpty(result)
	for index, rating := range result {
		rating.CreatedAt = time.Time{}
		rating.UpdatedAt = time.Time{}

		result[index] = rating
	}

	suite.Equal(model.Ratings{
		{
			Model:        gorm.Model{ID: 1},
			Score:        5,
			FoodRecipeID: 1,
			UserID:       "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
		},
		{
			Model:        gorm.Model{ID: 2},
			Score:        3,
			FoodRecipeID: 1,
			UserID:       "38fa4e9e-27de-42d5-a70f-9f01d41f32c2",
		},
	}, result)
}

func (suite *RepositoryGetRatingTestSuite) TestReturnEmptyWhenNotFound() {
	recipeID := 2

	// ไม่มี rating สำหรับ recipeID 2 ใน initial data
	result, err := suite.repository.Get(recipeID)

	suite.NoError(err)
	suite.Empty(result)
	suite.Equal(0, len(result))
}

func TestRepositoryGetRatings(t *testing.T) {
	suite.Run(t, new(RepositoryGetRatingTestSuite))
}
