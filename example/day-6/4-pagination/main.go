package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// slide: API Pagination

func main() {
	router := gin.Default()
	router.GET("/movies", func(ctx *gin.Context) {
		page := 1
		limit := 10

		if p := ctx.Query("page"); p != "" {
			if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
				page = parsed
			}
		}

		if ps := ctx.Query("limit"); ps != "" {
			if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 {
				limit = parsed
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"page":  page,
			"limit": limit,
		})
	})

	router.Run(":8000")
}
