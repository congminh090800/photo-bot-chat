package main

import (
	"os"
	"time"

	"github.com/congminh090800/photo-bot-chat/service"
	"github.com/go-co-op/gocron"
)

func startWorker() {
	scheduler := gocron.NewScheduler(time.UTC).SingletonMode()

	/**
	--------------- FETCH MEDIA ------------------
	*/
	photoService := service.NewPhotoService()
	scheduler.Every(55).Minutes().Do(func() {
		photoService.SaveMedia(os.Getenv("TARGET_ALBUM_ID"), "")
	})
	scheduler.StartAsync()
}
