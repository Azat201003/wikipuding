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

func GetByUserId(db *gorm.DB, id uint) ([]Article, error) {
	articles_list := []Article{}
	err := db.Find(&articles_list, &Article{CreatorId: id}).Error
	return articles_list, err
}

func GetById(db *gorm.DB, id uint) (*Article, error) {
	article := new(Article)
	err := db.Take(&article, Article{ID: id}).Error
	return article, err
}
