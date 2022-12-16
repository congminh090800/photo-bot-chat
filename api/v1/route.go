package v1_route

import (
	auth_route "github.com/congminh090800/photo-bot-chat/api/v1/auth"
	media_route "github.com/congminh090800/photo-bot-chat/api/v1/media"
	"github.com/gin-gonic/gin"
)

func AddRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	auth_route.AddRoutes(auth)

	media := rg.Group("/media")
	media_route.AddRoutes(media)
}
