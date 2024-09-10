package handlers

import (
	"errors"
	"log"

	"github.com/Ckefa/ckefablog.git/db"
	"github.com/Ckefa/ckefablog.git/models"
	"github.com/labstack/echo/v4"
)

func Subscribe(c echo.Context) error {
	new_email := c.FormValue("email")

	if new_email == "" {
		log.Println("Email not Found!!")
		return errors.New("Mail not Found")
	}

	if db.DB == nil {
		return errors.New("Database not Connected!")
	}

	user := models.NewUser(new_email)
	log.Println("crating user", user)
	db.DB.Create(user)

	return c.Render(200, "thanks", nil)
}

func Login(c echo.Context) error {
	return c.Render(200, "login", nil)
}

func Signup(c echo.Context) error {
	return c.Render(200, "signup", nil)
}
