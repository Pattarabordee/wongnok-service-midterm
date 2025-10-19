package helper

import (
	"net/http"
	"wongnok/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func DecodeClaims(ctx *gin.Context) (model.Claims, error) {
	value, exists := ctx.Get("claims")
	if !exists {
		return model.Claims{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	claims, ok := value.(model.Claims)
	if !ok {
		return model.Claims{}, errors.New(http.StatusText(http.StatusUnauthorized))
	}

	return claims, nil
}
