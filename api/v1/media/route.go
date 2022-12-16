package media_route

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(rg *gin.RouterGroup) {
	rg.GET("/", HandleMedia)
	rg.GET("/timelines/:date_type", HandleTimelines)
}
