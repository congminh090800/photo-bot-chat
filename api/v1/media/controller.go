package media_route

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/congminh090800/photo-bot-chat/model"
	"github.com/congminh090800/photo-bot-chat/service"
	"github.com/gin-gonic/gin"
)

func HandleTimelines(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": r.(error).Error(),
				"code":  http.StatusInternalServerError,
			})
			return
		}
	}()

	dateType, ok := ctx.Params.Get("date_type")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("no data type found").Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}
	photoService := service.NewPhotoService()
	data, err := photoService.GetTimelines(model.NewDateType(dateType))
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
		"code": http.StatusOK,
	})
}

func HandleMedia(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": r.(error).Error(),
				"code":  http.StatusInternalServerError,
			})
			return
		}
	}()

	dateType := ctx.Query("date_type")

	dateRange := ctx.Query("date_range")

	limit := ctx.Query("limit")
	intLimit, _ := strconv.Atoi(limit)

	offset := ctx.Query("offset")
	intOffset, _ := strconv.Atoi(offset)
	if intOffset < 0 {
		intOffset = 0
	}
	switch {
	case intLimit > 100:
		intLimit = 100
	case intLimit <= 0:
		intLimit = 10
	}
	photoService := service.NewPhotoService()
	data, err := photoService.GetMediaByTimeline(dateRange, model.NewDateType(dateType), intLimit, intOffset)
	if err != nil {
		panic(err)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
		"code": http.StatusOK,
	})
}
