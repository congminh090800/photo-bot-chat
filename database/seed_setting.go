package database

import (
	"log"

	"github.com/congminh090800/photo-bot-chat/model"
	"gorm.io/gorm/clause"
)

var setting = []model.Setting{
	{
		Key:   "MEDIA_NEXT_PAGE_TOKEN",
		Value: "",
	},
	{
		Key:   "GOOGLE_GLOBAL_ACCESS_TOKEN",
		Value: "",
	},
	{
		Key:   "GOOGLE_GLOBAL_REFRESH_TOKEN",
		Value: "",
	},
	{
		Key:   "GOOGLE_GLOBAL_TOKEN_TYPE",
		Value: "",
	},
	{
		Key:   "GOOGLE_GLOBAL_EXPIRY",
		Value: "",
	},
}

func SeedAllSetting() {
	db := GetDB()
	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(setting).Error
	if err != nil {
		log.Println("Failed to seed settings")
	}
	log.Println("Done seeding settings")
}
