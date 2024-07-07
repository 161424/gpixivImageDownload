package sql

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

//type DA struct {
//	*sql.PixivArtist
//}

type Artist struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
	ArtistId         string         `gorm:"column:ArtistId;type:text;primaryKey" json:"userId"`
	ArtistName       string         `gorm:"column:ArtistName;type:text" json:"name"`
	ArtistAvatar     string         `gorm:"column:ArtistAvatar;type:text" json:"image"`
	ArtistAvatarBig  string         `gorm:"column:ArtistAvatarBig;type:text" json:"imageBig"`
	ArtistBackground string         `gorm:"column:ArtistBackground;type:text" json:"background"`
	// 下载的输入参数
	BookMark    bool           `gorm:"column:BookMark"`
	Tags        pq.StringArray `gorm:"column:Tags"`
	Profile     bool           `gorm:"column:Profile"`
	DwImageInfo []Work         `gorm:"column:DwList"`
	DwImageUrl  []string
	Limit       int `gorm:"column:Limit"`
	Total       int
	InDb        bool `gorm:"column:InDb"`

	Illusts    pq.StringArray `gorm:"column:Illusts;type:text[]" json:"illusts"`
	Manga      pq.StringArray `gorm:"column:Manga;type:text[]" json:"manga"`
	Novels     pq.StringArray `gorm:"column:Novels;type:text[]" json:"novels"`
	IsLastPage bool           `gorm:"column:IsLastPage"`
	HaveImages bool           `gorm:"column:HaveImages"`
	Offset     int            `gorm:"column:Offset"`

	MangaSeries pq.StringArray `gorm:"column:MangaSeries;type:text[]" json:"mangaSeries"`
	NovelSeries pq.StringArray `gorm:"column:NovelSeries;type:text[]" json:"novelSeries"`
	Error       string         `json:"error"`
}

type AuthorWorks struct {
	Works []Work `json:"works"`
	Total int    `json:"total"`
	Error string `json:"error"`
}

type Work struct {
	Id         string   `json:"id"`
	Title      string   `json:"title"`
	IllustType int      `json:"illustType"`
	Tags       []string `json:"tags"`
	UserId     string   `json:"userId"`
	UserName   string   `json:"userName"`
	Width      float64  `json:"width"`
	Height     float64  `json:"height"`
	PageCount  int      `json:"pageCount"`
	CreateDate int      `json:"createDate"`
	AiType     int      `json:"aiType"`
}
