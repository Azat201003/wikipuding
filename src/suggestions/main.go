package suggestions

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Azat201003/wikipuding/src/auth"
	"github.com/Azat201003/wikipuding/src/likes"
	"github.com/Azat201003/wikipuding/src/wiki/articles"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SuggestionCreate struct {
	Token         string `header:"token"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	BaseArticleId uint   `param:"id"`
}

func getSuggestions(id uint, db *gorm.DB, ctx context.Context, client *redis.Client) string {
	result := "<ol>"
	suggestions := []articles.Article{}
	db.Model(&articles.Article{}).Find(&suggestions, &articles.Article{BaseArticleId: id})
	for _, suggestion := range suggestions {
		user := &auth.User{}
		db.Model(&auth.User{}).First(user, &auth.User{ID: suggestion.CreatorId})
		a, b := likes.CountLikes(suggestion.ID, ctx, client)
		result += fmt.Sprintf("<li>[<a href=\"../../%v/\">%v</a>] %v\t<i>by %v </i><u><b>%v</b> likes</u>\n%v</li>", suggestion.ID, suggestion.ID, suggestion.Title, user.Username, a-b, getSuggestions(suggestion.ID, db, ctx, client))
	}
	result += "</ol>"
	return result
}

func Init(e *echo.Echo, db *gorm.DB, ctx context.Context, client *redis.Client) {

	e.GET("/wiki/:id/suggestions/", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("GET /wiki/:id/suggestions/\terror with parsing id: %v\n", err.Error())
			return err
		}
		result := getSuggestions(uint(id), db, ctx, client)
		log.Printf("GET /wiki/:id/suggestions/\tarticle id: %v\n", uint(id))
		return c.HTML(http.StatusOK, result)
	})

	e.POST("/wiki/:id/suggestions/", func(c echo.Context) error {
		suggestion_create := new(SuggestionCreate)
		if err := c.Bind(suggestion_create); err != nil {
			log.Printf("POST /wiki/:id/suggestions/\terror with parsing data: %v\n", err.Error())
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, suggestion_create); err != nil {
			log.Printf("POST /wiki/:id/suggestions/\terror with parsing headers: %v\n", err.Error())
			return err
		}
		if err := b.BindPathParams(c, suggestion_create); err != nil {
			log.Printf("POST /wiki/:id/suggestions/\terror with parsing headers: %v\n", err.Error())
			return err
		}
		creator := new(auth.User)
		db.First(creator, &auth.User{Token: suggestion_create.Token})
		article := new(articles.Article)
		article.Title = suggestion_create.Title
		article.Content = suggestion_create.Content
		article.BaseArticleId = suggestion_create.BaseArticleId
		article.IsBase = false
		article.CreatorId = creator.ID
		db.Create(article)
		log.Printf("POST /wiki/:id/suggestions/\tarticle id: %v\n", article.ID)
		return c.JSON(http.StatusCreated, article.ID)
	})
}
