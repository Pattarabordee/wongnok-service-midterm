package auth_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"slices"
	"testing"
	"wongnok/internal/auth"
	"wongnok/internal/config"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

func TestNewHandler(t *testing.T) {
	t.Run("ShouldFillProperties", func(t *testing.T) {
		oauth2Conf := new(MockIOAuth2Config)
		verifier := new(MockIOIDCTokenVerifier)

		handler := auth.NewHandler(&gorm.DB{}, config.Keycloak{Realm: "demo"}, oauth2Conf, verifier)

		value := reflect.Indirect(reflect.ValueOf(handler))

		for index := 0; index < value.NumField(); index++ {
			field := value.Field(index)
			assert.False(t, field.IsZero(), "Field %s is zero value", field.Type().Name())
		}
	})
}

type HandlerTestSuite struct {
	suite.Suite

	// Dependencies
	handler     auth.IHandler
	service     *MockIService
	userService *MockIUserService
}

func (suite *HandlerLoginTestSuite) SetupSuite() {
	// Gin testing mode
	gin.SetMode(gin.TestMode)
}

// Extend
type HandlerLoginTestSuite struct {
	HandlerTestSuite

	// Mock data
	respGenerateState string
	respAuthCodeURL   string

	// Helper
	server func(payload io.Reader) *httptest.ResponseRecorder
}

func (suite *HandlerLoginTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.userService = new(MockIUserService)
	suite.handler = &auth.Handler{
		Service:     suite.service,
		UserService: suite.userService,
	}

	suite.server = func(payload io.Reader) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()
		router.GET("/api/v1/login", suite.handler.Login)

		// ใช้ httptest.NewRecorder() จำลอง server และ request เพื่อตรวจสอบ HTTP response จริงๆ
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodGet,
			"/api/v1/login",
			payload,
		)
		suite.NoError(err)

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respGenerateState = "state"
	suite.respAuthCodeURL = "http://example.com"

	suite.service.On("GenerateState").Return(func() string {
		return suite.respGenerateState
	})
	suite.service.On("AuthCodeURL", mock.Anything).Return(func(string) string {
		return suite.respAuthCodeURL
	})
}

func (suite *HandlerLoginTestSuite) TestRedirectToAuthURL() {
	response := suite.server(nil)

	body := response.Result().Body
	defer body.Close()

	suite.True(slices.ContainsFunc(response.Result().Cookies(), func(obj *http.Cookie) bool {
		return obj.Name == "state" && obj.Value == "state"
	}))

	suite.Equal(http.StatusTemporaryRedirect, response.Code)
	suite.Equal("http://example.com", response.Header().Get("Location"))

	suite.service.AssertCalled(suite.T(), "GenerateState")
	suite.service.AssertCalled(suite.T(), "AuthCodeURL", "state")
}

func TestHandlerLogin(t *testing.T) {
	suite.Run(t, new(HandlerLoginTestSuite))
}

// Extend
type HandlerCallbackTestSuite struct {
	HandlerTestSuite

	// Mock data
	respExchange         model.Credential
	errExchange          error
	respVerifyToken      *MockIOIDCIDToken
	errVerifyToken       error
	errClaims            error
	respUpsertWithClaims model.User
	errUpsertWithClaims  error

	// Helper
	server func(params url.Values) *httptest.ResponseRecorder
}

func (suite *HandlerCallbackTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.userService = new(MockIUserService)
	suite.handler = &auth.Handler{
		Service:     suite.service,
		UserService: suite.userService,
	}

	suite.server = func(params url.Values) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()
		router.GET("/api/v1/callback", suite.handler.Callback)

		// ใช้ httptest.NewRecorder() จำลอง server และ request เพื่อตรวจสอบ HTTP response จริงๆ
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodGet,
			"/api/v1/callback",
			nil,
		)
		suite.NoError(err)

		// Add query params
		query := request.URL.Query()
		for key := range params {
			query.Add(key, params.Get(key))
		}
		request.URL.RawQuery = query.Encode()

		// add cookie to request
		request.AddCookie(&http.Cookie{
			Name:     "state",
			Value:    "state",
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   300,
			HttpOnly: true,
			Raw:      "state=state; Path=/; Domain=localhost; Max-Age=300; HttpOnly",
		})

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respExchange = model.Credential{
		Token: &oauth2.Token{
			AccessToken: "token",
		},
		IDToken: "idtoken",
	}
	suite.errExchange = nil
	suite.respVerifyToken = new(MockIOIDCIDToken)
	suite.errVerifyToken = nil
	suite.errClaims = nil
	suite.respUpsertWithClaims = model.User{
		ID: "id",
	}
	suite.errUpsertWithClaims = nil

	suite.service.On("Exchange", mock.Anything, mock.Anything).Return(func(context.Context, string) (model.Credential, error) {
		return suite.respExchange, suite.errExchange
	})
	suite.service.On("VerifyToken", mock.Anything, mock.Anything).Return(func(ctx context.Context, token string) (auth.IOIDCIDToken, error) {
		return suite.respVerifyToken, suite.errVerifyToken
	})
	suite.respVerifyToken.On("Claims", mock.Anything).Run(func(args mock.Arguments) {
		claims := args.Get(0).(*model.Claims)
		*claims = model.Claims{
			ID: "id",
		}
	}).Return(func(any) error {
		return suite.errClaims
	})
	suite.userService.On("UpsertWithClaims", mock.Anything).Return(func(model.Claims) (model.User, error) {
		return suite.respUpsertWithClaims, suite.errUpsertWithClaims
	})
}

func (suite *HandlerCallbackTestSuite) TestResponseCredential() {
	params := url.Values{}
	params.Add("state", "state")
	params.Add("code", "code")

	response := suite.server(params)

	body := response.Result().Body
	defer body.Close()

	expectedResponse := dto.CredentialResponse{
		AccessToken: "token",
		IDToken:     "idtoken",
	}
	expectedJson, _ := json.Marshal(expectedResponse)

	suite.Equal(http.StatusOK, response.Code)
	suite.Equal(string(expectedJson), response.Body.String())

	suite.service.AssertCalled(suite.T(), "Exchange", mock.Anything, "code")
	suite.service.AssertCalled(suite.T(), "VerifyToken", mock.Anything, "idtoken")
	suite.userService.AssertCalled(suite.T(), "UpsertWithClaims", model.Claims{ID: "id"})
}

func (suite *HandlerCallbackTestSuite) TestErrorWhenCheckCookie() {
	params := url.Values{}
	params.Add("state", "1")

	response := suite.server(params)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusBadRequest, response.Code)
	suite.Equal(`{"message":"Invalid state"}`, response.Body.String())

	suite.service.AssertNotCalled(suite.T(), "Exchange", mock.Anything, mock.Anything)
}

func (suite *HandlerCallbackTestSuite) TestErrorWhenExchange() {
	suite.errExchange = assert.AnError

	params := url.Values{}
	params.Add("state", "state")
	params.Add("state", "code")

	response := suite.server(params)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())

	suite.service.AssertNotCalled(suite.T(), "VerifyToken", mock.Anything, mock.Anything)
}

func (suite *HandlerCallbackTestSuite) TestErrorWhenVerifyToken() {
	suite.errVerifyToken = assert.AnError

	params := url.Values{}
	params.Add("state", "state")
	params.Add("state", "code")

	response := suite.server(params)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())

	suite.userService.AssertNotCalled(suite.T(), "UpsertWithClaims", mock.Anything)
}

func (suite *HandlerCallbackTestSuite) TestErrorWhenDecodeClaims() {
	suite.errClaims = assert.AnError

	params := url.Values{}
	params.Add("state", "state")
	params.Add("state", "code")

	response := suite.server(params)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())

	suite.userService.AssertNotCalled(suite.T(), "UpsertWithClaims", mock.Anything)
}

func (suite *HandlerCallbackTestSuite) TestErrorWhenUpsertWithClaims() {
	suite.errUpsertWithClaims = assert.AnError

	params := url.Values{}
	params.Add("state", "state")
	params.Add("state", "code")

	response := suite.server(params)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())
}

func TestHandlerCallback(t *testing.T) {
	suite.Run(t, new(HandlerCallbackTestSuite))
}

// Extend
type HandlerLogoutTestSuite struct {
	HandlerTestSuite

	// Mock data
	respLogoutURL string
	errLogoutURL  error

	// Helper
	server func(params url.Values) *httptest.ResponseRecorder
}

func (suite *HandlerLogoutTestSuite) SetupTest() {
	suite.service = new(MockIService)
	suite.userService = new(MockIUserService)
	suite.handler = &auth.Handler{
		Service:     suite.service,
		UserService: suite.userService,
	}

	suite.server = func(params url.Values) *httptest.ResponseRecorder {
		// Create router
		router := gin.Default()
		router.GET("/api/v1/logout", suite.handler.Logout)

		// ใช้ httptest.NewRecorder() จำลอง server และ request เพื่อตรวจสอบ HTTP response จริงๆ
		recorder := httptest.NewRecorder()

		// Create request
		request, err := http.NewRequest(
			http.MethodGet,
			"/api/v1/logout",
			nil,
		)
		suite.NoError(err)

		// Add query params
		query := request.URL.Query()
		for key := range params {
			query.Add(key, params.Get(key))
		}
		request.URL.RawQuery = query.Encode()

		// Start testing server
		router.ServeHTTP(recorder, request)

		return recorder
	}

	suite.respLogoutURL = "http://example.com"
	suite.errLogoutURL = nil

	suite.service.On("LogoutURL", mock.Anything).Return(func(dto.LogoutQuery) (string, error) {
		return suite.respLogoutURL, suite.errLogoutURL
	})
}

func (suite *HandlerLogoutTestSuite) TestRedirectToLogoutURL() {
	params := url.Values{}
	params.Add("idTokenHint", "token")

	response := suite.server(params)

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusTemporaryRedirect, response.Code)
	suite.Equal("http://example.com", response.Header().Get("Location"))

	suite.service.AssertCalled(suite.T(), "LogoutURL", dto.LogoutQuery{IDTokenHint: "token"})
}

func (suite *HandlerLogoutTestSuite) TestErrorWhenGetLogoutURL() {
	suite.errLogoutURL = assert.AnError

	response := suite.server(url.Values{})

	body := response.Result().Body
	defer body.Close()

	suite.Equal(http.StatusInternalServerError, response.Code)
	suite.Equal(`{"message":"assert.AnError general error for testing"}`, response.Body.String())
}

func TestHandlerLogout(t *testing.T) {
	suite.Run(t, new(HandlerLogoutTestSuite))
}
