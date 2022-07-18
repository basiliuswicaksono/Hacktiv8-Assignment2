package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	Customer_Name string `gorm:"type:varchar(100)"`
	Ordered_At    time.Time
	Items         []Item
}
