package main

import (
	"net/http"
	_ "wongnok/example/day-11/docs"
	"wongnok/example/day-11/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Teacher API
// @version 1.0
// @description This is an example server.
// @host localhost:8000
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/students", getStudents)
	router.POST("/students", createStudent)

	router.Run(":8000")
}

// GetStudents godoc
// @Summary Get a student
// @Description get student
// @Tags students
// @Accept json
// @Produce json
// @Param page query int true "Page"
// @Param limit query int true "Limit"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /students [get]
func getStudents(ctx *gin.Context) {
	var query dto.PaginationQuery

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"page":  query.Page,
		"limit": query.Limit,
	})
}

// CreateStudents godoc
// @Summary Create a student
// @Description get student
// @Tags students
// @Accept json
// @Produce json
// @Param student body dto.StudentRequest true "Student Request"
// @Success 200 {object} dto.StudentResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /students [post]
func createStudent(ctx *gin.Context) {

	var request dto.StudentRequest

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	uuid := uuid.New().String()

	response := dto.StudentResponse{
		ID:        uuid,
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	ctx.JSON(http.StatusOK, response)
}
