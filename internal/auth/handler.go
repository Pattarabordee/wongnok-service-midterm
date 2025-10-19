package auth

import (
	"net/http"
	"wongnok/internal/config"
	"wongnok/internal/model"
	"wongnok/internal/model/dto"
	"wongnok/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IUserService user.IService

type IHandler interface {
	Login(ctx *gin.Context)
	Callback(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type Handler struct {
	Service     IService
	UserService IUserService
}

func NewHandler(db *gorm.DB, kc config.Keycloak, oauth2Conf IOAuth2Config, verifier IOIDCTokenVerifier) IHandler {
	return &Handler{
		Service:     NewService(kc, oauth2Conf, verifier),
		UserService: user.NewService(db),
	}
}

func (handler Handler) Login(ctx *gin.Context) {
	// Generate state
	state := handler.Service.GenerateState()

	// Collect state in Cookie
	ctx.SetCookie("state", state, 300, "/", "localhost", false, true)

	// Redirect to Keycloak
	ctx.Redirect(http.StatusTemporaryRedirect, handler.Service.AuthCodeURL(state))
}

func (handler Handler) Callback(ctx *gin.Context) {
	var query dto.KeycloakCallbackQuery
	ctx.BindQuery(&query)

	// Verify state
	state, err := ctx.Cookie("state")
	if err != nil || query.State != state {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid state"})
		return
	}

	// Exchange token
	credential, err := handler.Service.Exchange(ctx.Request.Context(), query.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Verify
	idToken, err := handler.Service.VerifyToken(ctx.Request.Context(), credential.IDToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Parse claims
	var claims model.Claims
	if err := idToken.Claims(&claims); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Ensure user
	if _, err := handler.UserService.UpsertWithClaims(claims); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, credential.ToResponse())
}

func (handler Handler) Logout(ctx *gin.Context) {
	var query dto.LogoutQuery
	ctx.BindQuery(&query)

	// Make logout url
	logoutURL, err := handler.Service.LogoutURL(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	// Redirect
	ctx.Redirect(http.StatusTemporaryRedirect, logoutURL)
}
