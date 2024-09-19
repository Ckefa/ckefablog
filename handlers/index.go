package handlers

import (
	"log"
	"net/http"

	"github.com/Ckefa/ckefablog/db"
	"github.com/Ckefa/ckefablog/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func HandleHome(c echo.Context) error {
	var cust models.Customer

	sess, _ := session.Get("session", c)
	user_id, ok := sess.Values["user_id"]
	if !ok {
		log.Println("User not loggged in!!")
	}

	if err := db.DB.Where("id = ?", user_id).First(&cust).Error; err != nil {
		log.Println(err)
	}
	user := map[string]interface{}{
		"name":  cust.Fname,
		"email": cust.Email,
		"subs":  cust.PackageID,
	}
	respData := map[string]interface{}{
		"user":  user,
		"packs": models.OrderDetails,
	}

	return c.Render(http.StatusOK, "home", respData)
}

func HandleIndex(c echo.Context) error {
	return c.Render(200, "index", nil)
}

func ServerSideRendering(c echo.Context) error {
	return c.Render(200, "introtossr", nil)
}

func GoogleFi(c echo.Context) error {
	return c.Render(200, "googlefi", nil)
}

func FutureOfRemoteWork(c echo.Context) error {
	return c.Render(200, "future_of_remote_work", nil)
}

func AiWorkplaceEthics(c echo.Context) error {
	return c.Render(200, "ai_workplace_ethics", nil)
}

func PrivacyPolicy(c echo.Context) error {
	return c.Render(200, "privacy_policy", nil)
}

func TermsOfService(c echo.Context) error {
	return c.Render(200, "terms_of_service", nil)
}

func About(c echo.Context) error {
	return c.Render(200, "about", nil)
}
