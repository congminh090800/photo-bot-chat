package api_route

import (
	v1_route "github.com/congminh090800/photo-bot-chat/api/v1"
	"github.com/gin-gonic/gin"
)

func AddRoutes(rg *gin.RouterGroup) {
	v1 := rg.Group("/v1")
	v1_route.AddRoutes(v1)
}
