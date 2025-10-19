package middleware

import (
	"net/http"
	"strings"
	"wongnok/internal/config"
	"wongnok/internal/model"

	"github.com/gin-gonic/gin"
)

func Authorize(verifier config.IOIDCTokenVerifier) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerPrefix := "Bearer "

		tokenWithBearer := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(tokenWithBearer, bearerPrefix) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": http.StatusText(http.StatusUnauthorized)})
			return
		}

		rawToken := strings.TrimPrefix(tokenWithBearer, bearerPrefix)
		idToken, err := verifier.Verify(ctx.Request.Context(), rawToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		var claims model.Claims
		if err := idToken.Claims(&claims); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		// Set claims in context
		ctx.Set("claims", claims)

		ctx.Next()
	}
}
