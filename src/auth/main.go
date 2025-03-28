package auth

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserSign struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

func generateToken() string {
	return "token3"
}

func GetByToken(db *gorm.DB, token string) (*User, error) {
	user := new(User)
	err := db.Take(user, &User{Token: token}).Error
	return user, err
}

func GetById(db *gorm.DB, id uint) (*User, error) {
	user := new(User)
	err := db.Take(&user, User{ID: id}).Error
	return user, err
}

func Init(e *echo.Echo, db *gorm.DB) {

	e.POST("/auth/sign-up/", func(c echo.Context) error {
		user_sign := new(UserSign)
		if err := c.Bind(user_sign); err != nil {
			log.Printf("POST /auth/sign-up/\terror with parsing data: %v\n", err.Error())
			return err
		}
		user := User{Username: user_sign.Username, Password: user_sign.Password, Level: 0, Token: generateToken()}
		db.Create(&user)
		log.Printf("POST /auth/sign-up/\tuser id: %v\n", user.ID)
		return c.JSON(http.StatusAccepted, user)
	})
	e.POST("/auth/sign-in/", func(c echo.Context) error {
		user_sign := new(UserSign)
		if err := c.Bind(user_sign); err != nil {
			log.Printf("POST /auth/sign-in/\terror with parsing data: %v\n", err.Error())
			return err
		}
		user := new(User)
		db.First(user, &User{Username: user_sign.Username, Password: user_sign.Password})
		log.Printf("POST /auth/sign-in/\tid: %v\n", user.ID)
		return c.JSON(http.StatusAccepted, user.Token)
	})
}
