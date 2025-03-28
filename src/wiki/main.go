package wiki

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Azat201003/wikipuding/src/auth"
	"github.com/Azat201003/wikipuding/src/likes"
	"github.com/Azat201003/wikipuding/src/wiki/articles"
)

func Init(e *echo.Echo, db *gorm.DB, ctx context.Context, client *redis.Client) {
	e.POST("/wiki/", func(c echo.Context) error {
		article_create := new(articles.ArticleCreate)
		if err := c.Bind(article_create); err != nil {
			log.Printf("POST /wiki/\terror with parsing data: %v\n", err.Error())
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, article_create); err != nil {
			log.Printf("POST /wiki/\terror with parsing headers: %v\n", err.Error())
			return err
		}
		creator := new(auth.User)
		err := db.First(creator, &auth.User{Token: article_create.Token}).Error
		if err != nil {
			log.Printf("POST /wiki/\terror with finding user: %v\n", err)
			return err
		}
		article := new(articles.Article)
		article.Title = article_create.Title
		article.Content = article_create.Content
		article.CreatorId = creator.ID
		article.IsBase = true
		err = db.Create(article).Error
		if err != nil {
			log.Printf("POST /wiki/\terror with creation article: %v\n", err)
			return err
		}
		log.Printf("POST /wiki/\tarticle id: %v\n", article.ID)
		return c.JSON(http.StatusCreated, article.ID)
	})

	e.POST("/wiki/", func(c echo.Context) error {
		article_create := new(articles.ArticleCreate)
		if err := c.Bind(article_create); err != nil {
			log.Printf("POST /wiki/\terror with parsing data: %v\n", err.Error())
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, article_create); err != nil {
			log.Printf("POST /wiki/\terror with parsing headers: %v\n", err.Error())
			return err
		}
		creator := new(auth.User)
		err := db.First(creator, &auth.User{Token: article_create.Token}).Error
		if err != nil {
			log.Printf("POST /wiki/\terror with finding user: %v\n", err)
			return err
		}
		article := new(articles.Article)
		article.Title = article_create.Title
		article.Content = article_create.Content
		article.CreatorId = creator.ID
		article.IsBase = true
		err = db.Create(article).Error
		if err != nil {
			log.Printf("POST /wiki/\terror with creation article: %v\n", err)
			return err
		}
		log.Printf("POST /wiki/\tarticle id: %v\n", article.ID)
		return c.JSON(http.StatusCreated, article.ID)
	})
	e.GET("/wiki/:id/", func(c echo.Context) error {
		var id int
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			log.Printf("GET /wiki/:id/\terror with parsing id: %v\n", err.Error())
			return err
		}

		article := new(articles.Article)
		err = db.First(&article, &articles.Article{ID: uint(id)}).Error
		if err != nil {
			log.Printf("GET /wiki/:id/\terror with getting article: %v\n", err.Error())
			return c.HTML(http.StatusNotFound, "<h1>article not found</h1>")
		}
		creator := new(auth.User)
		err = db.First(creator, auth.User{ID: article.CreatorId}).Error
		if err != nil {
			log.Printf("GET /wiki/:id/\terror with getting creator: %v\n", err.Error())
			return err
		}
		a, b := likes.CountLikes(article.ID, ctx, client)
		log.Printf("GET /wiki/:id/\taticle id: %v\n", article.ID)
		return c.HTML(http.StatusFound, fmt.Sprintf(`<h1>%v</h1><p>%v</p><i>by %v</i><br><b>%v</b> likes<br><b>%v</b> dislikes<br><a href="suggestions/">suggestions</a>`, article.Title, article.Content, creator.Username, a, b))
	})
	e.GET("/wiki/", func(c echo.Context) error {
		articles_list := []articles.Article{}
		db.Find(&articles_list, &articles.Article{IsBase: true})
		result := "<ul>"
		for _, article := range articles_list {
			user := new(auth.User)
			db.First(user, auth.User{ID: article.CreatorId})
			a, b := likes.CountLikes(article.ID, ctx, client)
			result += fmt.Sprintf(`<li>[<a href="%v/">%v</a>] <b>%v</b> <i>by %v</i> <u><b>%v</b> likes</u></li>`, article.ID, article.ID, article.Title, user.Username, a-b)
		}
		result += "</ul>"
		log.Printf("GET /wiki/\n")
		return c.HTML(http.StatusOK, result)
	})
}
