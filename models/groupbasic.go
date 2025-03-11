package models

import "gorm.io/gorm"

type GroupBasic struct {
	gorm.Model
	Name    string
	OnwerId uint
	Icon    string
	Notice  string
	Desc    string
}
