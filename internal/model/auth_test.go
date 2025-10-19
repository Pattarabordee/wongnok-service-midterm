package model_test

import (
	"testing"
	"time"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestCredentialToResponse(t *testing.T) {
	mockTime := time.Date(2025, 7, 19, 0, 0, 0, 0, time.Local)

	t.Run("ShouldTransformCredentialToResponse", func(t *testing.T) {
		cred := model.Credential{
			Token: &oauth2.Token{
				AccessToken:  "AccessToken",
				TokenType:    "TokenType",
				RefreshToken: "RefreshToken",
				Expiry:       mockTime,
				ExpiresIn:    3600,
			},
			IDToken: "IDToken",
		}

		expectedResponse := dto.CredentialResponse{
			AccessToken:  "AccessToken",
			TokenType:    "TokenType",
			RefreshToken: "RefreshToken",
			Expiry:       mockTime,
			ExpiresIn:    3600,
			IDToken:      "IDToken",
		}

		assert.Equal(t, expectedResponse, cred.ToResponse())
	})
}
