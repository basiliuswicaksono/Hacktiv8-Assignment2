package database

import (
	"assignment2/models"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	POSTGRES_HOST   = "172.17.0.3"
	POSTGRES_PORT   = "5432"
	POSTGRES_USER   = "postgres"
	POSTGRES_PASS   = "postgres"
	POSTGRES_DBNAME = "Assignment2-H8"
)

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST,
		POSTGRES_PORT,
		POSTGRES_USER,
		POSTGRES_PASS,
		POSTGRES_DBNAME,
	)

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	db.Debug().AutoMigrate(models.Order{}, models.Item{})
	return db
}
