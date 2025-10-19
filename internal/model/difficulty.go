package model

import "gorm.io/gorm"

type Difficulty struct {
	gorm.Model
	Name string
}
