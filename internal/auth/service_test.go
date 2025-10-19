package auth_test

import (
	"context"
	"encoding/base64"
	"reflect"
	"strings"
	"testing"
	"wongnok/internal/auth"
	"wongnok/internal/config"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/coreos/go-oidc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
)

func TestNewService(t *testing.T) {
	t.Run("ShouldFillProperties", func(t *testing.T) {
		oauth2Conf := new(MockIOAuth2Config)
		verifier := new(MockIOIDCTokenVerifier)

		service := auth.NewService(config.Keycloak{Realm: "demo"}, oauth2Conf, verifier)

		value := reflect.Indirect(reflect.ValueOf(service))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})
}

type ServiceTestSuite struct {
	suite.Suite

	// Dependencies
	keycloak     config.Keycloak
	oauth2config *MockIOAuth2Config
	verifier     *MockIOIDCTokenVerifier
	service      auth.IService
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.keycloak = config.Keycloak{
		Realm: "demo",
		URL:   "http://example.com",
	}
	suite.oauth2config = new(MockIOAuth2Config)
	suite.verifier = new(MockIOIDCTokenVerifier)
	suite.service = &auth.Service{
		Keycloak:     suite.keycloak,
		OAuth2Config: suite.oauth2config,
		Verifier:     suite.verifier,
	}
}

// Extend
type ServiceGenerateStateTestSuite struct {
	ServiceTestSuite
}

func (suite *ServiceGenerateStateTestSuite) TestReturnState() {
	state := suite.service.GenerateState()
	decoded, err := base64.URLEncoding.DecodeString(state)
	suite.NoError(err)

	suite.Len(decoded, 32)
}

func TestServiceGenerateState(t *testing.T) {
	suite.Run(t, new(ServiceGenerateStateTestSuite))
}

// Extend
type ServiceAuthCodeURLTestSuite struct {
	ServiceTestSuite

	// Mock data
	respAuthCodeURL string
}

func (suite *ServiceAuthCodeURLTestSuite) SetupTest() {
	// Super
	suite.ServiceTestSuite.SetupTest()

	suite.respAuthCodeURL = "http://example.com"

	suite.oauth2config.On("AuthCodeURL", mock.Anything).Return(func(string, ...oauth2.AuthCodeOption) string {
		return suite.respAuthCodeURL
	})
}

func (suite *ServiceAuthCodeURLTestSuite) TestReturnAuthCodeURL() {
	uri := suite.service.AuthCodeURL("state")
	suite.Equal("http://example.com", uri)
}

func TestServiceAuthCodeURL(t *testing.T) {
	suite.Run(t, new(ServiceAuthCodeURLTestSuite))
}

// Extend
type ServiceExchangeTestSuite struct {
	ServiceTestSuite

	// Mock data
	respExchange *oauth2.Token
	errExchange  error
}

func (suite *ServiceExchangeTestSuite) SetupTest() {
	// Super
	suite.ServiceTestSuite.SetupTest()

	token := &oauth2.Token{}
	suite.respExchange = token.WithExtra(map[string]any{"id_token": "token"})
	suite.errExchange = nil

	suite.oauth2config.On("Exchange", mock.Anything, mock.Anything).Return(func(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
		return suite.respExchange, suite.errExchange
	})
}

func (suite *ServiceExchangeTestSuite) TestReturnCredential() {
	expectedCred := model.Credential{
		Token:   suite.respExchange,
		IDToken: "token",
	}

	cred, err := suite.service.Exchange(context.Background(), "code")
	suite.NoError(err)

	suite.Equal(expectedCred, cred)
}

func (suite *ServiceExchangeTestSuite) TestErrorWhenExchange() {
	suite.errExchange = assert.AnError

	cred, err := suite.service.Exchange(context.Background(), "code")
	suite.ErrorIs(err, assert.AnError)
	suite.True(strings.HasPrefix(err.Error(), "exchange token"))

	suite.Empty(cred)
}

func (suite *ServiceExchangeTestSuite) TestErrorWhenGetExtra() {
	suite.respExchange = suite.respExchange.WithExtra(map[string]any{"id_token": nil})

	cred, err := suite.service.Exchange(context.Background(), "code")
	suite.EqualError(err, "id token is missing")

	suite.Empty(cred)
}

func TestServiceExchange(t *testing.T) {
	suite.Run(t, new(ServiceExchangeTestSuite))
}

// Extend
type ServiceVerifyTokenTestSuite struct {
	ServiceTestSuite

	// Mock data
	respVerify *oidc.IDToken
	errVerify  error
}

func (suite *ServiceVerifyTokenTestSuite) SetupTest() {
	// Super
	suite.ServiceTestSuite.SetupTest()

	suite.respVerify = &oidc.IDToken{Issuer: "demo"}
	suite.errVerify = nil

	suite.verifier.On("Verify", mock.Anything, mock.Anything).Return(func(context.Context, string) (*oidc.IDToken, error) {
		return suite.respVerify, suite.errVerify
	})
}

func (suite *ServiceVerifyTokenTestSuite) TestReturnToken() {
	token, err := suite.service.VerifyToken(context.Background(), "token")
	suite.NoError(err)

	suite.Equal(&oidc.IDToken{Issuer: "demo"}, token)
}

func (suite *ServiceVerifyTokenTestSuite) TestErrorWhenVerify() {
	suite.errVerify = assert.AnError

	token, err := suite.service.VerifyToken(context.Background(), "token")
	suite.ErrorIs(err, assert.AnError)

	suite.Empty(token)
}

func TestServiceVerifyToken(t *testing.T) {
	suite.Run(t, new(ServiceVerifyTokenTestSuite))
}

// Extend
type ServiceLogoutURLTestSuite struct {
	ServiceTestSuite
}

func (suite *ServiceLogoutURLTestSuite) TestReturnLogoutURL() {
	query := dto.LogoutQuery{
		IDTokenHint:           "token",
		PostLogoutRedirectURI: "redirect",
	}

	uri, err := suite.service.LogoutURL(query)
	suite.NoError(err)

	suite.Equal("http://example.com/realms/demo/protocol/openid-connect/logout?id_token_hint=token&post_logout_redirect_uri=redirect", uri)
}

func (suite *ServiceLogoutURLTestSuite) TestErrorWhenParseURL() {
	suite.keycloak.URL = ":/example.com/ path"
	suite.service = &auth.Service{
		Keycloak:     suite.keycloak,
		OAuth2Config: suite.oauth2config,
		Verifier:     suite.verifier,
	}

	uri, err := suite.service.LogoutURL(dto.LogoutQuery{})
	suite.Error(err)
	suite.True(strings.HasPrefix(err.Error(), "parse logout url"))

	suite.Empty(uri)
}

func TestServiceLogoutURL(t *testing.T) {
	suite.Run(t, new(ServiceLogoutURLTestSuite))
}
