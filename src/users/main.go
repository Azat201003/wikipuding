package users

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

func Init(e *echo.Echo, db *gorm.DB, ctx context.Context, client *redis.Client) {

	e.GET("/users/:id/", func(c echo.Context) error {
		var id int
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("GET /users/:id/\terror with parsing id: %v\n", err.Error())
			return err
		}
		user, err := auth.GetById(db, uint(id))
		if err != nil {
			log.Printf("GET /users/:id/\terror error with finding user by id user id: %v: %v\n", id, err.Error())
			return err
		}
		articles_list, err := articles.GetByUserId(db, uint(id))
		if err != nil {
			log.Printf("GET /users/:id/\terror with finding articles by user id user id: %v: %v\n", id, err.Error())
			return err
		}
		result := fmt.Sprintf(`<h1>%v</h1>Articles count: <i>%v</i>.<h5>articles: </h3>`, user.Username, len(articles_list))
		for _, article := range articles_list {
			a, b := likes.CountLikes(article.ID, ctx, client)
			result += fmt.Sprintf(`<li>[<a href="../../wiki/%v/">%v</a>] <b>%v</b> <i>by %v</i> <u><b>%v</b> likes</u></li>`, article.ID, article.ID, article.Title, user.Username, a-b)
		}
		log.Printf("GET /users/:id/\tuser id: %v\n", user.ID)
		return c.HTML(http.StatusFound, result)
	})

	e.GET("/users/", func(c echo.Context) error {
		users := []auth.User{}
		db.Find(&users)
		result := "<ul>"
		for _, user := range users {
			articles_list, err := articles.GetByUserId(db, user.ID)
			if err != nil {
				log.Printf("GET /users/\terror user id: %v: %v\n", user.ID, err.Error())
				return err
			}
			result += fmt.Sprintf(`<li>[<a href="%v/">%v</a>] <b>%v</b> <i>states count: %v</i></li>`, user.ID, user.ID, user.Username, len(articles_list))
		}
		result += "</ul>"
		log.Printf("GET\t/users/\n")
		return c.HTML(http.StatusOK, result)
	})

}
