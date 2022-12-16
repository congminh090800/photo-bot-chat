package auth_route

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/congminh090800/photo-bot-chat/database"
	"github.com/congminh090800/photo-bot-chat/model"
	"github.com/congminh090800/photo-bot-chat/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func HandleSocialLogin(ctx *gin.Context) {
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
	socialType, ok := ctx.Params.Get("social_type")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errors.New("no data type found").Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}
	socialType = strings.Trim(strings.ToLower(socialType), " ")
	handler := service.SocialLoginFactory(socialType)
	url := handler.CreateLoginUrl()
	// redirect to google sign in
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func HandleCallback(ctx *gin.Context) {
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

	// get social login handler
	socialType := ctx.Query("state")
	handler := service.SocialLoginFactory(socialType)

	// exchange code
	config := service.GetConfig()
	code := ctx.Query("code")
	token, err := config.Exchange(ctx, code)
	if err != nil {
		panic(err)
	}
	data, err := handler.GetUserInfo(token.AccessToken)
	if err != nil {
		panic(err)
	}

	session := sessions.Default(ctx)
	session.Set("user", data)
	session.Save()
	if data.Email == "congminh090800@gmail.com" {
		db := database.GetDB()
		err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&model.Setting{}).Where("key = ?", "GOOGLE_GLOBAL_ACCESS_TOKEN").Update("value", token.AccessToken).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Setting{}).Where("key = ?", "GOOGLE_GLOBAL_REFRESH_TOKEN").Update("value", token.RefreshToken).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Setting{}).Where("key = ?", "GOOGLE_GLOBAL_TOKEN_TYPE").Update("value", token.TokenType).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.Setting{}).Where("key = ?", "GOOGLE_GLOBAL_EXPIRY").Update("value", token.Expiry.Format(time.UnixDate)).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"data": data,
		"code": http.StatusOK,
	})
}
