package helper_test

import (
	"net/http/httptest"
	"testing"
	"wongnok/internal/helper"
	"wongnok/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDecodeClaims(t *testing.T) {
	t.Run("ShouldReturnClaims", func(t *testing.T) {
		expectedClaims := model.Claims{
			ID: "ID",
		}

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("claims", expectedClaims)

		claims, err := helper.DecodeClaims(ctx)
		assert.NoError(t, err)

		assert.Equal(t, expectedClaims, claims)
	})

	t.Run("ShouldBeErrorWhenGetClaimsFromContext", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

		claims, err := helper.DecodeClaims(ctx)
		assert.EqualError(t, err, "Unauthorized")

		assert.Empty(t, claims)
	})

	t.Run("ShouldBeErrorWhenDecodeClaimsToModel", func(t *testing.T) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("claims", "demo")

		claims, err := helper.DecodeClaims(ctx)
		assert.EqualError(t, err, "Unauthorized")

		assert.Empty(t, claims)
	})
}
