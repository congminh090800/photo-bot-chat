package service

import (
	"errors"
	"io"
	"net/http"
	"os"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var authLock = sync.Once{}
var googleConfig *oauth2.Config

func GetConfig() *oauth2.Config {
	authLock.Do(func() {
		googleConfig = &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("API_URL") + "/api/v1/auth/google/callback",
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"https://www.googleapis.com/auth/photoslibrary.readonly",
				"https://www.googleapis.com/auth/userinfo.profile",
				"https://www.googleapis.com/auth/userinfo.email",
			},
		}
	})
	return googleConfig
}

type SocialLoginHandler interface {
	CreateLoginUrl() string
	GetUserInfo(accessToken string) (*UserInfoSchema, error)
}

type UserInfoSchema struct {
	ID     string `json_google:"id"`
	Email  string `json_google:"email"`
	Name   string `json_google:"name"`
	Avatar string `json_google:"picture"`
}

/*
*
-------------- GOOGLE -------------------
*/
type GoogleLoginHandler struct {
	UserInfoApi string
}

func (g GoogleLoginHandler) CreateLoginUrl() string {
	conf := GetConfig()
	url := conf.AuthCodeURL(
		"google",
		oauth2.SetAuthURLParam("access_type", "offline"),
		oauth2.SetAuthURLParam("prompt", "consent"),
	)
	return url
}

func (g GoogleLoginHandler) GetUserInfo(accessToken string) (*UserInfoSchema, error) {
	resp, err := http.Get(g.UserInfoApi + "?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(contents))
	}
	userInfo := UserInfoSchema{}
	json := jsoniter.Config{
		TagKey: "json_google",
	}.Froze()
	err = json.Unmarshal(contents, &userInfo)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

/**
-------------- GOOGLE -------------------
*/

func SocialLoginFactory(socialType string) SocialLoginHandler {
	switch socialType {
	case "google":
		return GoogleLoginHandler{UserInfoApi: "https://www.googleapis.com/oauth2/v2/userinfo"}
	default:
		panic(errors.New("socialType not supported"))
	}
}
