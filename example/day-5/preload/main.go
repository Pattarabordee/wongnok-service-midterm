package main

import (
	"encoding/json"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Difficulty struct {
	gorm.Model
	Name string
}

type CookingDuration struct {
	gorm.Model
	Name string
}

type FoodRecipe struct {
	gorm.Model
	Name              string
	CookingDurationID uint
	CookingDuration   CookingDuration
	DifficultyID      uint
	Difficulty        Difficulty
}

func main() {
	db, err := gorm.Open(postgres.Open("postgres://postgres:pass2word@localhost:5432/wongnok"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	var recipe FoodRecipe

	if err := db.Preload("Difficulty").Preload("CookingDuration").First(&recipe).Error; err != nil {
		log.Fatal(err)
	}

	disp, _ := json.MarshalIndent(recipe, "", "  ")
	fmt.Println(string(disp))
}
