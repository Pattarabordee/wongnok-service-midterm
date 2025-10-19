package model

import (
	"wongnok/internal/model/dto"

	"golang.org/x/oauth2"
)

type Credential struct {
	*oauth2.Token
	IDToken string
}

func (cred Credential) ToResponse() dto.CredentialResponse {
	return dto.CredentialResponse{
		AccessToken:  cred.AccessToken,
		TokenType:    cred.TokenType,
		RefreshToken: cred.RefreshToken,
		Expiry:       cred.Expiry,
		ExpiresIn:    cred.ExpiresIn,
		IDToken:      cred.IDToken,
	}
}

type Claims struct {
	ID        string `json:"sub" validate:"required"`
	FirstName string `json:"given_name" validate:"required"`
	LastName  string `json:"family_name" validate:"required"`
	NickName  string `json:"nick_name"`
}
