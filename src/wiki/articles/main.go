package articles

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	ID            uint
	Title         string
	Content       string
	CreatorId     uint
	BaseArticleId uint
	IsBase        bool
}

type ArticleCreate struct {
	Token   string `header:"token"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
