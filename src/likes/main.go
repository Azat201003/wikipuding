package likes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Azat201003/wikipuding/src/auth"
	"github.com/Azat201003/wikipuding/src/wiki/articles"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type LikePost struct {
	Token  string `header:"token"`
	IsLike bool   `json:"is_like"`
}

func CountLikes(id uint, ctx context.Context, client *redis.Client) (int, int) {
	a := 0
	b := 0
	likes := client.SInter(ctx, fmt.Sprintf(`likes:%v`, id)).Val()
	for _, like := range likes {
		if like[0] != '-' {
			a++
		} else {
			b++
		}
	}

	return a, b
}

func Init(e *echo.Echo, db *gorm.DB, ctx context.Context, client *redis.Client) {

	e.POST("/wiki/:id/like/", func(c echo.Context) error {
		like_post := LikePost{}
		if err := c.Bind(&like_post); err != nil {
			log.Printf("POST /wiki/:id/like/\terror with parsing data: %v\n", err.Error())
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, &like_post); err != nil {
			log.Printf("POST /wiki/:id/like/\terror with parsing headers: %v\n", err.Error())
			return err
		}
		article := new(articles.Article)
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("POST /wiki/:id/like/\terror with parsing id: %v\n", err.Error())
			return err
		}
		err = db.First(article, &articles.Article{ID: uint(id)}).Error
		if err != nil {
			log.Printf("POST /wiki/:id/like/\terror with finding article article id: %v: %v\n", id, err)
			return err
		}
		count_likes := -1
		user := auth.User{}
		err = db.Take(&user, &auth.User{Token: like_post.Token}).Error
		if like_post.Token == "" {
			log.Printf("POST /wiki/:id/like/\terror with finding user: empty token\n")
			return err
		}
		if err != nil {
			log.Printf("POST /wiki/:id/like/\terror with finding user token: %v: %v\n", like_post.Token, err)
			return err
		}
		user_id := user.ID
		if like_post.IsLike {
			count_likes = 1
		}
		if client.SIsMember(ctx, fmt.Sprintf("likes:%v", article.ID), fmt.Sprintf("%v", int(user_id)*(-count_likes))).Val() {
			err = client.SRem(ctx, fmt.Sprintf("likes:%v", article.ID), fmt.Sprintf("%v", int(user_id)*(-count_likes))).Err()
			if err != nil {
				log.Printf("POST /wiki/:id/like/\terror with removing from set user id: %v, article id: %v, count likes: %v\n: %v", user.ID, article.ID, int(user_id)*(-count_likes), err)
				return err
			}
			log.Printf("POST /wiki/:id/like/\tarticle id: %v user id: %v like reseted\n", article.ID, user.ID)
			return c.JSON(http.StatusOK, "reseted")
		} else if client.SIsMember(ctx, fmt.Sprintf("likes:%v", article.ID), fmt.Sprintf("%v", int(user_id)*(count_likes))).Val() {
			log.Printf("POST /wiki/:id/like/\tarticle id: %v user id: %v like already was added\n", article.ID, user.ID)
			return c.JSON(http.StatusOK, "already liked")
		} else {
			err := client.SAdd(ctx, fmt.Sprintf("likes:%v", article.ID), []string{fmt.Sprintf("%v", int(user_id)*(count_likes))}).Err()
			if err != nil {
				log.Printf("POST /wiki/:id/like/\terror with adding like into set user id: %v, article id: %v, count likes: %v\n: %v", user.ID, article.ID, int(user_id)*(count_likes), err)
				return err
			}
			log.Printf("POST /wiki/:id/like/\tarticle id: %v user id: %v like added\n", article.ID, user.ID)
			return c.JSON(http.StatusOK, "liked")
		}
	})
}
