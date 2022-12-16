package model

type Album struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	ProductUrl      string `json:"productUrl"`
	MediaItemsCount int64  `json:"mediaItemsCount,string"`
}

type AlbumListResponse struct {
	SharedAlbums  []Album `json:"sharedAlbums"`
	NextPageToken string  `json:"nextPageToken"`
}
