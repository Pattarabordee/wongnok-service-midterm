//////////////////// [ Day 5 ] ////////////////////

package user_test

import (
	"testing"
	"wongnok/example/day-5/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServiceGet(t *testing.T) {
	t.Run("ShouldBeReturnUser", func(t *testing.T) {
		mockRepo := new(MockIRepository)
		service := user.NewService(mockRepo)

		mockRepo.On("Get", mock.Anything).Return(user.User{Name: "Peter"}, nil)

		result, err := service.Get("id")
		assert.NoError(t, err)

		assert.Equal(t, user.User{Name: "Peter"}, result)
	})

	t.Run("ShouldBeError", func(t *testing.T) {
		mockRepo := new(MockIRepository)
		service := user.NewService(mockRepo)

		mockRepo.On("Get", mock.Anything).Return(user.User{}, assert.AnError)

		result, err := service.Get("id")
		assert.ErrorIs(t, err, assert.AnError)

		assert.Empty(t, result)
	})
}

//////////////////// [ Day 5 ] ////////////////////
