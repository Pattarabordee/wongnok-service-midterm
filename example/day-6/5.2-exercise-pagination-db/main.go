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

// slide: exercise Pagination with GORM

type Teacher struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Age  int
}

func paginateTeachers(db *gorm.DB, page, limit int) ([]Teacher, int64, error) {
	var teachers []Teacher
	var total int64

	offset := (page - 1) * limit

	// Fetch teachers with pagination
	result := db.
		Limit(limit).
		Offset(offset).
		Find(&teachers)

	// Count total records
	db.Model(&Teacher{}).Count(&total)

	return teachers, total, result.Error
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

	router.GET("/teachers", func(ctx *gin.Context) {
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

		teachers, total, err := paginateTeachers(db, page, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching teachers"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"results": teachers,
			"total":   total,
		})
	})

	router.Run(":8000")
}
