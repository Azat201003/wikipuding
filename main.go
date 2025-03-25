package main

import (
	"fmt"
	"log"
	"strconv"

	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
)

type Article struct {
	gorm.Model
	ID        uint
	Title     string
	Content   string
	CreatorId uint
}

type User struct {
	gorm.Model
	ID       uint
	Level    uint
	Username string
	Token    string
	Password string
}
type ArticleCreate struct {
	Token   string `header:"token"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UserSign struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func generateToken() string {
	return "token3"
}

func main() {
	dsn := "host=localhost user=wiki password=1234 dbname=wiki port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	err = db.AutoMigrate(&Article{})
	if err != nil {
		log.Panic(err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Panic(err)
	}
	// db.Create(&Article{Title: "Article", Content: "It is very cool article!", CreatorId: 1})
	// var article Article
	// db.First(&article, &Article{Title: "Article"})
	// fmt.Println(article)
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "<h1>Hello, world!</h1>")
	})
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
	e.POST("/wiki/", func(c echo.Context) error {
		article_create := new(ArticleCreate)
		fmt.Println("Hello, world!")
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
		return c.HTML(http.StatusFound, fmt.Sprintf(`<h1>%v</h1><p>%v</p><i>%v<i>`, article.Title, article.Content, creator.Username))
	})
	e.GET("/wiki/", func(c echo.Context) error {
		articles := []Article{}
		db.Find(&articles)
		result := "<ul>"
		for _, article := range articles {
			user := new(User)
			db.First(user, User{ID: article.CreatorId})
			result += fmt.Sprintf(`<li>[<a href="%v/">%v</a>] <b>%v</b> <i>by %v</i></li>`, article.ID, article.ID, article.Title, user.Username)
		}
		result += "</ul>"
		return c.HTML(http.StatusOK, result)
	})
	e.Logger.Fatal(e.Start(":1323"))
}
