package main

import (
	"fmt"
	"log"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	ID            uint
	Title         string
	Content       string
	CreatorId     uint
	BaseArticleId uint
	IsBase        bool
}

type User struct {
	gorm.Model
	ID       uint
	Level    uint
	Username string
	Token    string
	Password string
	IsStaff  bool
}
type ArticleCreate struct {
	Token   string `header:"token"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SuggestionCreate struct {
	Token         string `header:"token"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	BaseArticleId uint   `json:"base_article_id"`
}

type UserSign struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func generateToken() string {
	return "token3"
}

func gerSuggestions(id uint, db *gorm.DB) string {
	result := "<ol>"
	suggestions := []Article{}
	db.Model(&Article{}).Find(&suggestions, &Article{BaseArticleId: id})
	for _, suggestion := range suggestions {
		user := &User{}
		db.Model(&User{}).First(user, &User{ID: suggestion.CreatorId})
		result += fmt.Sprintf("<li>[<a href=\"../../%v/\">%v</a>] %v\t<i>by %v</i>\n%v</li>", suggestion.ID, suggestion.ID, suggestion.Title, user.Username, gerSuggestions(suggestion.ID, db))
	}
	result += "</ol>"
	return result
}

func main() {
	dsn := "host=localhost user=wiki password=1234 dbname=wiki port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	// err = db.AutoMigrate(&Article{})
	// if err != nil {
	// 	log.Panic(err)
	// }
	// err = db.AutoMigrate(&User{})
	// if err != nil {
	// 	log.Panic(err)
	// }
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Hello!</h1><p>You can read some <a href=\"wiki/\">articles</a> or look for some <a href=\"users/\">user profiles</a></p>")
	})

	// * auth

	e.POST("/auth/sign-up/", func(c echo.Context) error {
		user_sign := new(UserSign)
		if err := c.Bind(user_sign); err != nil {
			return err
		}
		user := User{Username: user_sign.Username, Password: user_sign.Password, Level: 0, Token: generateToken()}
		db.Create(&user)
		return c.JSON(http.StatusAccepted, user)
	})
	e.POST("/auth/sign-in/", func(c echo.Context) error {
		user_sign := new(UserSign)
		if err := c.Bind(user_sign); err != nil {
			return err
		}
		user := new(User)
		db.First(user, &User{Username: user_sign.Username, Password: user_sign.Password})
		return c.JSON(http.StatusAccepted, user.Token)
	})

	// * suggestions

	e.GET("/wiki/:id/suggestions/", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		result := gerSuggestions(uint(id), db)
		return c.HTML(http.StatusOK, result)
	})
	e.POST("/wiki/:id/suggestions/", func(c echo.Context) error {
		suggestion_create := new(SuggestionCreate)
		if err := c.Bind(suggestion_create); err != nil {
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, suggestion_create); err != nil {
			return err
		}
		creator := new(User)
		db.First(creator, &User{Token: suggestion_create.Token})
		article := new(Article)
		article.Title = suggestion_create.Title
		article.Content = suggestion_create.Content
		article.BaseArticleId = suggestion_create.BaseArticleId
		article.IsBase = false
		article.CreatorId = creator.ID
		db.Create(article)
		return c.JSON(http.StatusCreated, article.ID)
	})

	// * articles (wiki)

	e.POST("/wiki/", func(c echo.Context) error {
		article_create := new(ArticleCreate)
		if err := c.Bind(article_create); err != nil {
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, article_create); err != nil {
			return err
		}
		creator := new(User)
		db.First(creator, &User{Token: article_create.Token})
		article := new(Article)
		article.Title = article_create.Title
		article.Content = article_create.Content
		article.CreatorId = creator.ID
		article.IsBase = true
		db.Create(article)
		return c.JSON(http.StatusCreated, article.ID)
	})
	e.GET("/wiki/:id/", func(c echo.Context) error {
		var id int
		id, _ = strconv.Atoi(c.Param("id"))
		article := new(Article)
		db.First(&article, &Article{ID: uint(id)})
		creator := new(User)
		db.First(creator, User{ID: article.CreatorId})
		return c.HTML(http.StatusFound, fmt.Sprintf(`<h1>%v</h1><p>%v</p><i>%v<i><br><a href="suggestions/">suggestions</a>`, article.Title, article.Content, creator.Username))
	})
	e.GET("/wiki/", func(c echo.Context) error {
		articles := []Article{}
		db.Find(&articles, &Article{IsBase: true})
		result := "<ul>"
		for _, article := range articles {
			user := new(User)
			db.First(user, User{ID: article.CreatorId})
			result += fmt.Sprintf(`<li>[<a href="%v/">%v</a>] <b>%v</b> <i>by %v</i></li>`, article.ID, article.ID, article.Title, user.Username)
		}
		result += "</ul>"
		return c.HTML(http.StatusOK, result)
	})

	// * users
	e.GET("/users/:id/", func(c echo.Context) error {
		var id int
		id, _ = strconv.Atoi(c.Param("id"))
		user := new(User)
		db.First(&user, &User{ID: uint(id)})
		articles := []Article{}
		db.Find(&articles, &Article{CreatorId: uint(id)})
		return c.HTML(http.StatusFound, fmt.Sprintf(`<h1>%v</h1>States count: <i>%v<i>.`, user.Username, len(articles)))
	})
	e.GET("/users/", func(c echo.Context) error {
		users := []User{}
		db.Find(&users)
		result := "<ul>"
		for _, user := range users {
			articles := []Article{}
			db.Find(&articles, &Article{CreatorId: user.ID})
			result += fmt.Sprintf(`<li>[<a href="%v/">%v</a>] <b>%v</b> <i>states count: %v</i></li>`, user.ID, user.ID, user.Username, len(articles))
		}
		result += "</ul>"
		return c.HTML(http.StatusOK, result)
	})
	e.Logger.Fatal(e.Start(":1323"))

}
