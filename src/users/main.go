package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Azat201003/wikipuding/src/auth"
	"github.com/Azat201003/wikipuding/src/wiki/articles"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Init(e *echo.Echo, db *gorm.DB) {

	e.GET("/users/:id/", func(c echo.Context) error {
		var id int
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Printf("GET /users/:id/\terror with parsing id: %v\n", err.Error())
			return err
		}
		user := new(auth.User)
		db.First(&user, &auth.User{ID: uint(id)})
		articles_list := []articles.Article{}
		db.Find(&articles_list, &articles.Article{CreatorId: uint(id)})
		log.Printf("GET /wiki/\tuser id: %v\n", user.ID)
		return c.HTML(http.StatusFound, fmt.Sprintf(`<h1>%v</h1>States count: <i>%v<i>.`, user.Username, len(articles_list)))
	})

	e.GET("/users/", func(c echo.Context) error {
		users := []auth.User{}
		db.Find(&users)
		result := "<ul>"
		for _, user := range users {
			articles_list := []articles.Article{}
			db.Find(&articles_list, &articles.Article{CreatorId: user.ID})
			result += fmt.Sprintf(`<li>[<a href="%v/">%v</a>] <b>%v</b> <i>states count: %v</i></li>`, user.ID, user.ID, user.Username, len(articles_list))
		}
		result += "</ul>"
		log.Printf("GET\t/users/\n")
		return c.HTML(http.StatusOK, result)
	})

}
