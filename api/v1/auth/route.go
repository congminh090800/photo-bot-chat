package auth_route

import (
	"github.com/gin-gonic/gin"
)

func AddRoutes(rg *gin.RouterGroup) {
	rg.GET("/:social_type", HandleSocialLogin)
	rg.GET("/:social_type/callback", HandleCallback)
}
