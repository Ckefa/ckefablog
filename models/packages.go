package models

import "gorm.io/gorm"

type Package struct {
	gorm.Model
	ID    int64   `gorm:"primaryKey" json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var Packages = []Package{
	{ID: 1, Name: "none", Price: 0},
	{ID: 2, Name: "basic", Price: 199},
	{ID: 3, Name: "standard", Price: 499},
	{ID: 4, Name: "premium", Price: 999},
}
