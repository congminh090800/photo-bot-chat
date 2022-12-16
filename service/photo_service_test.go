package service

import (
	"log"
	"testing"

	"github.com/congminh090800/photo-bot-chat/database"
	"github.com/congminh090800/photo-bot-chat/model"
	"github.com/joho/godotenv"
)

var TARGET_ALBUM_ID = "ABgK9mx2llS5IHRGLOvxdYm1-8ndpOmfsg1o99wtTL4X4LVY80deM2fY6jCKqKoaX1qvWTjc6EUYMyXIQ8Btwb10Xoes7u4pIg"

func setupTest(t *testing.T) func(t *testing.T) {
	log.Println("Initializing environment variables...")
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatalf("Error loading .env file")
	}
	log.Println("Initializing database connection...")
	database.ConnectDB()
	// Return a function to teardown the test
	return func(t *testing.T) {
		db := database.GetDB()
		if sqlDB, err := db.DB(); err != nil {
			t.Errorf("Failed to close database")
		} else {
			sqlDB.Close()
			log.Println("Closed database connection!")
		}
	}
}

func TestDbRelatedCases(t *testing.T) {
	tearDownTest := setupTest(t)
	defer tearDownTest(t)

	t.Run("TestGetMediaTimeline", func(t *testing.T) {
		service := NewPhotoService()
		timelines, err := service.GetTimelines(model.Month)
		if err != nil {
			t.Errorf(err.Error())
		}
		if len(timelines) == 0 {
			t.Errorf("must atleast 1 timeline")
		}
		for _, p := range timelines {
			log.Printf("%+v\n", p)
		}
	})

	// t.Run("TestGetSharedAlbums", func(t *testing.T) {
	// 	service := NewPhotoService()
	// 	albumListResponse, err := service.GetSharedAlbums("", "")
	// 	if err != nil {
	// 		t.Errorf(err.Error())
	// 	}
	// 	if len(albumListResponse.SharedAlbums) == 0 {
	// 		t.Errorf("must atleast 1 album")
	// 	}
	// })

	// t.Run("TestGetMediaList", func(t *testing.T) {
	// 	pService := NewPhotoService()
	// 	mediaList, err := pService.GetMediaList(TARGET_ALBUM_ID, "", "")
	// 	if err != nil {
	// 		t.Errorf(err.Error())
	// 	}
	// 	if len(mediaList.MediaItems) == 0 {
	// 		t.Errorf("must atleast 1 media")
	// 	}
	// })

	// t.Run("TestBulkCreateMedia", func(t *testing.T) {
	// 	pService := NewPhotoService()
	// 	ok := pService.SaveMedia(TARGET_ALBUM_ID, "")
	// 	if !ok {
	// 		t.Errorf("failed to fetch media")
	// 	}
	// })

}
