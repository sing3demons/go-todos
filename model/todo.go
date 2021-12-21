package model

import "gorm.io/gorm"

type Todo struct {
	gorm.Model
	Title string `gorm:"not null"`
	Desc  string
	Image string
}
