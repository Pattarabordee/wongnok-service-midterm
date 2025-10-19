package model

import "gorm.io/gorm"

type CookingDuration struct {
	gorm.Model
	Name string
}
