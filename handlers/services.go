package handlers

import (
	"log"
	"net/http"

	"github.com/Ckefa/ckefablog/db"
	"github.com/Ckefa/ckefablog/models"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func GetServices(c echo.Context) error {
	var cust models.Customer
	var pack models.Package

	sess, _ := session.Get("session", c)

	userID := sess.Values["user_id"]
	if err := db.DB.Where("id = ?", userID).First(&cust).Error; err != nil {
		log.Println("Func: GetServices - Error getting customer", err)
	}

	if err := db.DB.Where("id = ?", cust.PackageID).First(&pack).Error; err != nil {
		log.Println("Func: GetServices - Error getting package ", err)
	}

	user := map[string]interface{}{
		"name":  cust.Fname,
		"email": cust.Email,
		"subs":  cust.PackageID,
	}

	respData := map[string]interface{}{
		"user":     user,
		"names":    cust.Fname + " " + cust.Lname,
		"pid":      cust.PackageID,
		"pack":     pack.Name,
		"price":    pack.Price,
		"status":   "inprogress",
		"progress": 10,
	}

	return c.Render(http.StatusOK, "services", respData)
}
