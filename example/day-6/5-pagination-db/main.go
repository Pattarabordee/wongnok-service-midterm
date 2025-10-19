package main

import (
	"log"
	"net/http"
	"strconv"
	"wongnok/internal/config"

	"github.com/caarlos0/env/v11"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// slide: Connect SQL to GORM Pagination

type Movie struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	SerialNumber string
}

func paginateMovies(db *gorm.DB, page, limit int) ([]Movie, int64, error) {
	var movies []Movie
	var total int64

	offset := (page - 1) * limit

	// Fetch movies with pagination
	result := db.
		Limit(limit).
		Offset(offset).
		Find(&movies)

	// Count total records
	db.Model(&Movie{}).Count(&total)

	return movies, total, result.Error
}

func main() {

	// Load configuration
	var conf config.Config

	if err := env.Parse(&conf); err != nil {
		log.Fatal("Error when decoding configuration:", err)
	}

	// Database connection
	db, err := gorm.Open(postgres.Open(conf.Database.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error when connect to database:", err)
	}

	// Ensure close connection when terminated
	defer func() {
		sqldb, _ := db.DB()
		sqldb.Close()
	}()

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

		movies, total, err := paginateMovies(db, page, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching movies"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"results": movies,
			"total":   total,
		})
	})

	router.Run(":8000")
}
