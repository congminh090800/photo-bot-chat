package main

import (
	"log"
	"net/http"
	"time"

	route_api "github.com/congminh090800/photo-bot-chat/api"
	"github.com/congminh090800/photo-bot-chat/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func main() {
	InitEnv()
	database.ConnectDB()
	startWorker()
	router.SetTrustedProxies([]string{"*"})
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowedMethods:   []string{"PUT", "PATCH", "GET", "OPTIONS", "POST", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		MaxAge:           12 * time.Hour,
	}))
	addSession()
	getRoutes()
	router.Run()
}

func getRoutes() {
	router.GET("/health", func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": r.(error).Error(),
					"code":  http.StatusInternalServerError,
				})
				return
			}
		}()
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	api := router.Group("/api")
	route_api.AddRoutes(api)

	router.Use(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, "not found")
	})
}

func addSession() {
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("photo-bot-chat", store))
}
