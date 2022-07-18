package models

import "github.com/jinzhu/gorm"

type Item struct {
	gorm.Model
	Item_Code   string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:varchar(500)"`
	Quantity    int
	OrderID     uint
}
