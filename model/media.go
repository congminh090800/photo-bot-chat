package model

import "time"

type MediaMetadata struct {
	CreationTimeStr string    `json:"creationTime" gorm:"-:all"` // eg. 2014-10-02T15:01:23Z
	CreationTime    time.Time `gorm:"not null"`
	Width           int64     `json:"width,string"`
	Height          int64     `json:"height,string"`
}

type Media struct {
	ID          string        `json:"id" gorm:"primaryKey"`
	Description string        `json:"description"`
	ProductUrl  string        `json:"productUrl"`
	BaseUrl     string        `json:"baseUrl" gorm:"not null"`
	MimeType    string        `json:"mimeType" gorm:"not null"`
	Filename    string        `json:"filename"`
	Metadata    MediaMetadata `json:"mediaMetadata" gorm:"embedded"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type MediaListResponse struct {
	MediaItems    []Media `json:"mediaItems"`
	NextPageToken string  `json:"nextPageToken"`
}

type MediaTimeline struct {
	DateRange time.Time
	Total     int
}
