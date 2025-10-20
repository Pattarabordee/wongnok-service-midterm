package model

import "time"

type RecipeLove struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"column:user_id"`
	RecipeID  int       `gorm:"column:recipe_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}
