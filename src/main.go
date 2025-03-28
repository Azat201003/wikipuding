package main

import (
	"context"
	"log"
	"os"

	"net/http"

	"github.com/Azat201003/wikipuding/src/auth"
	"github.com/Azat201003/wikipuding/src/likes"
	"github.com/Azat201003/wikipuding/src/suggestions"
	"github.com/Azat201003/wikipuding/src/users"
	"github.com/Azat201003/wikipuding/src/wiki"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
		log.Print("Failed to open log file: ", err)
	}

	log.SetOutput(file)

	// * home page

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		log.Printf("GET /\n")
		return c.HTML(http.StatusOK, "<h1>Hello!</h1><p>You can read some <a href=\"wiki/\">articles</a> or look for some <a href=\"users/\">user profiles</a></p>")
	})

	// * auth
	auth.Init(e, db)

	// * suggestions
	suggestions.Init(e, db, ctx, client)

	// * articles (wiki)
	wiki.Init(e, db, ctx, client)

	// * users
	users.Init(e, db, ctx, client)

	// * likes
	likes.Init(e, db, ctx, client)
	e.Logger.Fatal(e.Start(":1323"))

}

/*
------------------------------------------------------------------------------
File                                       blank        comment           code
------------------------------------------------------------------------------
./src/wiki/main.go                             7              0            113
./src/likes/main.go                            7              0             91
./src/suggestions/main.go                      7              0             72
./docker-compose.yaml                          8              2             47
./src/auth/main.go                             7              0             47
./main.go                                     16              6             43
./src/users/main.go                            6              0             40
./src/wiki/articles/main.go                    3              0             16
./ideas.yaml                                   0              0             10
./Dockerfile                                   4              3              6
./main_test.go                                 0              0              1
-------------------------------------------------------------------------------
SUM:                                          65             11            486
-------------------------------------------------------------------------------

-------------------------------------------------------------------------------
Language                     files          blank        comment           code
-------------------------------------------------------------------------------
Go                               8             53              6            423
YAML                             2              8              2             57
Dockerfile                       1              4              3              6
-------------------------------------------------------------------------------
SUM:                            11             65             11            486
-------------------------------------------------------------------------------
*/
