package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/congminh090800/photo-bot-chat/database"
	"github.com/congminh090800/photo-bot-chat/model"
	"golang.org/x/oauth2"
	"gorm.io/gorm/clause"
)

const GOOGLE_PHOTO_LIBRARY_API_URL = "https://photoslibrary.googleapis.com/v1"

type PhotoService struct {
	Token *oauth2.Token
	Http  *http.Client
}

func NewPhotoService() *PhotoService {
	db := database.GetDB()
	var accessTokenSetting, refreshTokenSetting, tokenTypeSetting, expirySetting model.Setting
	if err := db.Where("key = ?", "GOOGLE_GLOBAL_ACCESS_TOKEN").First(&accessTokenSetting).Error; err != nil {
		panic(err)
	}
	if err := db.Where("key = ?", "GOOGLE_GLOBAL_REFRESH_TOKEN").First(&refreshTokenSetting).Error; err != nil {
		panic(err)
	}
	if err := db.Where("key = ?", "GOOGLE_GLOBAL_TOKEN_TYPE").First(&tokenTypeSetting).Error; err != nil {
		panic(err)
	}
	if err := db.Where("key = ?", "GOOGLE_GLOBAL_EXPIRY").First(&expirySetting).Error; err != nil {
		panic(err)
	}
	expiry, err := time.Parse(time.UnixDate, expirySetting.Value)
	if err != nil {
		panic(err)
	}

	instance := &PhotoService{
		Token: &oauth2.Token{
			AccessToken:  accessTokenSetting.Value,
			RefreshToken: refreshTokenSetting.Value,
			TokenType:    tokenTypeSetting.Value,
			Expiry:       expiry,
		},
	}
	return instance
}

func NewPhotoServiceFromToken(token *oauth2.Token) *PhotoService {
	instance := &PhotoService{
		Token: token,
	}
	return instance
}

func (p *PhotoService) GetSharedAlbums(pageToken string, accessToken string) (*model.AlbumListResponse, error) {
	var httpClient *http.Client
	apiUrl := GOOGLE_PHOTO_LIBRARY_API_URL + "/sharedAlbums?pageSize=50"
	if pageToken != "" {
		apiUrl = apiUrl + "&pageToken=" + pageToken
	}

	if accessToken == "" {
		httpClient = GetConfig().Client(context.Background(), p.Token)
	} else {
		httpClient = &http.Client{}
		apiUrl = apiUrl + "&access_token=" + accessToken
	}

	resp, err := httpClient.Get(apiUrl)
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

	result := model.AlbumListResponse{}
	err = json.Unmarshal(contents, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *PhotoService) GetAlbum(id string, accessToken string) (*model.Album, error) {
	pageToken := ""
	for {
		albumListResponse, err := p.GetSharedAlbums(pageToken, accessToken)
		if err != nil {
			return nil, err
		}
		pageToken = albumListResponse.NextPageToken
		for _, item := range albumListResponse.SharedAlbums {
			if item.ID == id {
				result := &model.Album{
					ID:              item.ID,
					Title:           item.Title,
					ProductUrl:      item.ProductUrl,
					MediaItemsCount: item.MediaItemsCount,
				}
				return result, nil
			}
		}
		if pageToken == "" {
			return nil, errors.New("end of list but album not found")
		}
	}
}

func (p *PhotoService) GetMediaList(albumId string, pageToken string, accessToken string) (*model.MediaListResponse, error) {
	var httpClient *http.Client
	if accessToken == "" {
		httpClient = GetConfig().Client(context.Background(), p.Token)
	} else {
		httpClient = &http.Client{}
	}

	apiUrl := GOOGLE_PHOTO_LIBRARY_API_URL + "/mediaItems:search"
	payload := map[string]any{
		"albumId":  albumId,
		"pageSize": 100,
		// "filters": map[string]any{
		// 	"dateFilter": map[string]any{
		// 		"ranges": []map[string]any{
		// 			{
		// 				"startDate": map[string]int{
		// 					"year":  fromDate.Year(),
		// 					"month": int(fromDate.Month()),
		// 					"day":   fromDate.Day(),
		// 				},
		// 				"endDate": map[string]int{
		// 					"year":  toDate.Year(),
		// 					"month": int(toDate.Month()),
		// 					"day":   toDate.Day(),
		// 				},
		// 			},
		// 		},
		// 	},
		// },
	}
	if pageToken != "" {
		payload["pageToken"] = pageToken
	}
	postBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}
	if accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+accessToken)
	}

	resp, err := httpClient.Do(req)
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

	result := model.MediaListResponse{}
	err = json.Unmarshal(contents, &result)
	if err != nil {
		return nil, err
	}
	for idx, item := range result.MediaItems {
		newTime, err := time.Parse(time.RFC3339, strings.TrimSuffix(item.Metadata.CreationTimeStr, "Z")+"+00:00")
		if err != nil {
			return nil, err
		}
		result.MediaItems[idx].Metadata.CreationTime = newTime
	}
	return &result, nil
}

func (p *PhotoService) SaveMedia(albumId string, accessToken string) bool {
	log.Printf("job:fetch-media:start")
	db := database.GetDB()
	pageToken := ""
	for {
		mediaListResponse, err := p.GetMediaList(albumId, pageToken, accessToken)
		if err != nil {
			log.Println("job:fetch-media:error-on-request", err)
			return false
		}
		log.Printf("job:fetch-media:found %d media", len(mediaListResponse.MediaItems))

		if len(mediaListResponse.MediaItems) > 0 {
			if err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&mediaListResponse.MediaItems).Error; err != nil {
				log.Println("job:fetch-media:error-on-upsert-media", err)
				return false
			}
		}

		pageToken = mediaListResponse.NextPageToken
		if pageToken == "" {
			break
		}
	}
	log.Printf("job:fetch-media:end")
	return true
}

func (p *PhotoService) GetTimelines(dateType model.DateType) ([]*model.MediaTimeline, error) {
	if dateType.String() == "" {
		return nil, errors.New("dateType not found")
	}

	db := database.GetDB()
	var results []*model.MediaTimeline
	err := db.Table("media").Select("date_trunc(?, media.creation_time at time zone 'utc') as date_range, count(*) as total", dateType.String()).Group("date_range").Order("date_range desc").Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (p *PhotoService) GetMediaByTimeline(dateRange string, dateType model.DateType, limit int, offset int) (*database.Paging, error) {
	if dateType.String() == "" {
		return nil, errors.New("dateType not found")
	}

	db := database.GetDB()
	var err error

	var total int64
	err = db.Model(&model.Media{}).Where("date_trunc(?, media.creation_time at time zone 'utc') =  ?", dateType.String(), dateRange).Count(&total).Error
	if err != nil {
		return nil, err
	}
	var results []*model.Media
	err = db.Scopes(database.Paginate(limit, offset)).Where("date_trunc(?, media.creation_time at time zone 'utc') =  ?", dateType.String(), dateRange).Find(&results).Error
	if err != nil {
		return nil, err
	}

	return &database.Paging{
		Items:  results,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}, nil
}
