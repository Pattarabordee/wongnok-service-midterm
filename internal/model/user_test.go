package model_test

import (
	"testing"
	"time"
	"wongnok/internal/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserFromClaims(t *testing.T) {
	mockTime := time.Date(2025, 7, 19, 0, 0, 0, 0, time.Local)

	t.Run("ShouldTransformClaimsToUser", func(t *testing.T) {
		claims := model.Claims{
			ID:        "ID",
			FirstName: "FirstName",
			LastName:  "LastName",
		}

		user := model.User{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
			},
		}

		expectedUser := model.User{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: mockTime,
				UpdatedAt: mockTime,
			},
			ID:        "ID",
			FirstName: "FirstName",
			LastName:  "LastName",
		}

		assert.Equal(t, expectedUser, user.FromClaims(claims))
	})
}
