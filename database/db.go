package database

import (
	"fmt"
	"log"
	"os"

	"github.com/congminh090800/photo-bot-chat/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	if DB == nil {
		panic("Please connect to database first")
	}
	return DB
}

func ConnectDB() {
	var err error
	config := map[string]string{
		"DBHost": os.Getenv("PGHOST"),
		"DBName": os.Getenv("PGDATABASE"),
		"DBUser": os.Getenv("PGUSER"),
		"DBPass": os.Getenv("PGPASSWORD"),
		"DBPort": os.Getenv("PGPORT"),
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Shanghai",
		config["DBHost"],
		config["DBUser"],
		config["DBPass"],
		config["DBName"],
		config["DBPort"],
	)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	DB.AutoMigrate(&model.Media{}, &model.Setting{})
	log.Println("Connected Successfully to the Database")
	SeedAllSetting()
}

func Paginate(limit int, offset int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		_offset := offset
		_limit := limit
		if offset < 0 {
			_offset = 0
		}
		switch {
		case limit > 100:
			_limit = 100
		case limit <= 0:
			_limit = 10
		}

		return db.Offset(_offset).Limit(_limit)
	}
}

type Paging struct {
	Items  any
	Total  int64
	Limit  int
	Offset int
}
