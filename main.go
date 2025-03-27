package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
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

type LikePost struct {
	Token  string `header:"token"`
	IsLike bool   `json:"is_like"`
}

func generateToken() string {
	return "token3"
}

func gerSuggestions(id uint, db *gorm.DB, ctx context.Context, client *redis.Client) string {
	result := "<ol>"
	suggestions := []Article{}
	db.Model(&Article{}).Find(&suggestions, &Article{BaseArticleId: id})
	for _, suggestion := range suggestions {
		user := &User{}
		db.Model(&User{}).First(user, &User{ID: suggestion.CreatorId})
		a, b := countLikes(suggestion.ID, ctx, client)
		result += fmt.Sprintf("<li>[<a href=\"../../%v/\">%v</a>] %v\t<i>by %v </i><u><b>%v</b> likes</u>\n%v</li>", suggestion.ID, suggestion.ID, suggestion.Title, user.Username, a-b, gerSuggestions(suggestion.ID, db, ctx, client))
	}
	result += "</ol>"
	return result
}

func countLikes(id uint, ctx context.Context, client *redis.Client) (int, int) {
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

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "1234",
		DB:       0,
	})
	ctx := context.Background()
	dsn := "host=localhost user=wiki password=1234 dbname=wiki port=5432 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	file, err := os.OpenFile("main.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}

	log.SetOutput(file)

	// * home page

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		log.Println("GET /")
		return c.HTML(http.StatusOK, "<h1>Hello!</h1><p>You can read some <a href=\"wiki/\">articles</a> or look for some <a href=\"users/\">user profiles</a></p>")
	})

	// * auth

	e.POST("/auth/sign-up/", func(c echo.Context) error {
		user_sign := new(UserSign)
		if err := c.Bind(user_sign); err != nil {
			log.Fatal(fmt.Sprintf("POST /auth/sign-up/ error with parsing data: %v\n", err.Error()))
			return err
		}
		user := User{Username: user_sign.Username, Password: user_sign.Password, Level: 0, Token: generateToken()}
		db.Create(&user)
		log.Printf("POST /auth/sign-up/ user id: %v\n", user.ID)
		return c.JSON(http.StatusAccepted, user)
	})
	e.POST("/auth/sign-in/", func(c echo.Context) error {
		user_sign := new(UserSign)
		if err := c.Bind(user_sign); err != nil {
			log.Fatal(fmt.Sprintf("POST /auth/sign-in/ error with parsing data: %v\n", err.Error()))
			return err
		}
		user := new(User)
		db.First(user, &User{Username: user_sign.Username, Password: user_sign.Password})
		log.Printf("POST /auth/sign-in/ id: %v\n", user.ID)
		return c.JSON(http.StatusAccepted, user.Token)
	})

	// * suggestions

	e.GET("/wiki/:id/suggestions/", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatal(fmt.Sprintf("GET /wiki/:id/suggestions/ error with parsing id: %v\n", err.Error()))
			return err
		}
		result := gerSuggestions(uint(id), db, ctx, client)
		log.Printf("GET /wiki/:id/suggestions/ article id: %v\n", uint(id))
		return c.HTML(http.StatusOK, result)
	})

	e.POST("/wiki/:id/suggestions/", func(c echo.Context) error {
		suggestion_create := new(SuggestionCreate)
		if err := c.Bind(suggestion_create); err != nil {
			log.Fatal(fmt.Sprintf("POST /wiki/:id/suggestions/ error with parsing data: %v\n", err.Error()))
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, suggestion_create); err != nil {
			log.Fatal(fmt.Sprintf("POST /wiki/:id/suggestions/ error with parsing headers: %v\n", err.Error()))
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
		log.Printf("POST /wiki/:id/suggestions/ article id: %v\n", article.ID)
		return c.JSON(http.StatusCreated, article.ID)
	})

	// * articles (wiki)

	e.POST("/wiki/", func(c echo.Context) error {
		article_create := new(ArticleCreate)
		if err := c.Bind(article_create); err != nil {
			log.Fatal(fmt.Sprintf("POST /wiki/ error with parsing data: %v\n", err.Error()))
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, article_create); err != nil {
			log.Fatal(fmt.Sprintf("POST /wiki/ error with parsing headers: %v\n", err.Error()))
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
		log.Printf("POST /wiki/ article id: %v\n", article.ID)
		return c.JSON(http.StatusCreated, article.ID)
	})
	e.GET("/wiki/:id/", func(c echo.Context) error {
		var id int
		id, err = strconv.Atoi(c.Param("id"))

		if err != nil {
			log.Fatalf("GET /wiki/:id/ error with parsing id: %v\n", err.Error())
			return err
		}

		article := new(Article)
		db.First(&article, &Article{ID: uint(id)})
		creator := new(User)
		db.First(creator, User{ID: article.CreatorId})
		a, b := countLikes(article.ID, ctx, client)
		log.Printf("GET /wiki/:id/ aticle id: %v", article.ID)
		return c.HTML(http.StatusFound, fmt.Sprintf(`<h1>%v</h1><p>%v</p><i>by %v</i><br><b>%v</b> likes<br><b>%v</b> dislikes<br><a href="suggestions/">suggestions</a>`, article.Title, article.Content, creator.Username, a, b))
	})
	e.GET("/wiki/", func(c echo.Context) error {
		articles := []Article{}
		db.Find(&articles, &Article{IsBase: true})
		result := "<ul>"
		for _, article := range articles {
			user := new(User)
			db.First(user, User{ID: article.CreatorId})
			a, b := countLikes(article.ID, ctx, client)
			result += fmt.Sprintf(`<li>[<a href="%v/">%v</a>] <b>%v</b> <i>by %v</i> <u><b>%v</b> likes</u></li>`, article.ID, article.ID, article.Title, user.Username, a-b)
		}
		result += "</ul>"
		log.Printf("GET /wiki/\n")
		return c.HTML(http.StatusOK, result)
	})

	// * users

	e.GET("/users/:id/", func(c echo.Context) error {
		var id int
		id, err = strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatalf("GET /users/:id/ error with parsing id: %v\n", err.Error())
			return err
		}
		user := new(User)
		db.First(&user, &User{ID: uint(id)})
		articles := []Article{}
		db.Find(&articles, &Article{CreatorId: uint(id)})
		log.Printf("GET /wiki/ user id: %v\n", user.ID)
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
		log.Printf("GET /users/\n")
		return c.HTML(http.StatusOK, result)
	})

	// * likes

	e.POST("/wiki/:id/like/", func(c echo.Context) error {
		like_post := LikePost{}
		if err := c.Bind(&like_post); err != nil {
			log.Fatal(fmt.Sprintf("POST /wiki/:id/like/ error with parsing data: %v\n", err.Error()))
			return err
		}
		b := &echo.DefaultBinder{}
		if err := b.BindHeaders(c, &like_post); err != nil {
			log.Fatal(fmt.Sprintf("POST /wiki/:id/like/ error with parsing headers: %v\n", err.Error()))
			return err
		}
		article := new(Article)
		id, _ := strconv.Atoi(c.Param("id"))
		db.First(article, &Article{ID: uint(id)})
		count_likes := -1
		user := new(User)
		db.First(user, &User{Token: like_post.Token})
		user_id := user.ID
		if like_post.IsLike {
			count_likes = 1
		}
		if client.SIsMember(ctx, fmt.Sprintf("likes:%v", article.ID), fmt.Sprintf("%v", int(user_id)*(-count_likes))).Val() {
			client.SRem(ctx, fmt.Sprintf("likes:%v", article.ID), fmt.Sprintf("%v", int(user_id)*(-count_likes)))
			log.Printf("POST /wiki/:id/like/ article id: %v user id: %v like reseted\n", article.ID, user.ID)
			return c.JSON(http.StatusOK, "reseted")
		} else if client.SIsMember(ctx, fmt.Sprintf("likes:%v", article.ID), fmt.Sprintf("%v", int(user_id)*(count_likes))).Val() {
			log.Printf("POST /wiki/:id/like/ article id: %v user id: %v like already was added\n", article.ID, user.ID)
			return c.JSON(http.StatusOK, "already liked")
		} else {
			client.SAdd(ctx, fmt.Sprintf("likes:%v", article.ID), []string{fmt.Sprintf("%v", int(user_id)*(count_likes))})
			log.Printf("POST /wiki/:id/like/ article id: %v user id: %v like added\n", article.ID, user.ID)
			return c.JSON(http.StatusOK, "liked")
		}
	})
	e.Logger.Fatal(e.Start(":1323"))

}
