package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Ckefa/ckefablog/db"
	"github.com/Ckefa/ckefablog/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Subscribe(c echo.Context) error {
	new_email := c.FormValue("email")

	if new_email == "" {
		log.Println("<< func Subscribe: - Email not Found!!")
		return errors.New("Mail not Found")
	}

	if db.DB == nil {
		return errors.New("<< func Subscribe: - Database not Connected!")
	}

	user := models.NewUser(new_email)
	log.Println("crating user", user)
	db.DB.Create(user)

	return c.Render(200, "thanks", nil)
}

func Register(c echo.Context) error {
	new_email := c.FormValue("email")
	fname := c.FormValue("fname")
	lname := c.FormValue("lname")
	passwd := c.FormValue("passwd")

	if new_email == "" {
		log.Println("Email not Found!!")
		return errors.New("Mail not Found")
	}

	if db.DB == nil {
		return errors.New("Database not Connected!")
	}

	cust := models.NewCustomer(fname, lname, new_email, passwd)
	log.Println("crating customer", cust)
	db.DB.Create(cust)

	c.Response().Header().Set("HX-Redirect", "/login")
	return c.NoContent(http.StatusOK)
}

func Login(c echo.Context) error {
	email := c.FormValue("email")
	passwd := c.FormValue("passwd")

	var cust models.Customer

	if err := db.DB.Where("email = ?", email).First(&cust).Error; err != nil {
		log.Println(err)
		return c.String(http.StatusAccepted, err.Error())
	}

	if passwd != cust.Passwd {
		return c.String(http.StatusAccepted, "Incorrect password")
	}

	sess, _ := session.Get("session", c)
	sess.Values["user_id"] = cust.ID
	sess.Save(c.Request(), c.Response())

	c.Response().Header().Set("HX-Redirect", "/services")
	return c.NoContent(http.StatusOK)
}

func HandleLogin(c echo.Context) error {
	return c.Render(200, "login", nil)
}

func Logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	respData := map[string]interface{}{
		"packs": models.OrderDetails,
	}

	return c.Render(http.StatusOK, "home", respData)

}

func Signup(c echo.Context) error {
	return c.Render(200, "signup", nil)
}
