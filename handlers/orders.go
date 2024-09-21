package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Ckefa/ckefablog/db"
	"github.com/Ckefa/ckefablog/models"
	"github.com/Ckefa/ckefablog/paypal"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Checkout(c echo.Context) error {
	pid := c.Param("pid")

	var pack models.Package

	if err := db.DB.Where("id = ?", pid).First(&pack).Error; err != nil {
		return c.JSON(http.StatusAccepted, map[string]string{"message": "no package selected"})
	}
	log.Println("Checking out package ", pack.Name)

	packData := map[string]interface{}{
		"name":    pack.Name,
		"price":   pack.Price,
		"details": models.OrderDetails[int(pack.ID)],
	}

	return c.Render(http.StatusOK, "checkout", packData)
}

func RequestOrder(c echo.Context) error {
	sess, _ := session.Get("session", c)
	userID, _ := sess.Values["user_id"].(string)
	amt := c.FormValue("amount")

	amt_float, err := strconv.ParseFloat(amt, 64)
	pid := models.GetPid(amt_float)

	log.Println("Paying $", amt_float)

	order := models.NewOrder(string(userID), pid, amt)
	err = paypal.CreateOrder(order)
	if err != nil {
		return c.String(http.StatusAccepted, "error making payment")
	}
	if err := db.DB.Create(&order).Error; err != nil {
		log.Println(err)
	}

	c.Response().Header().Set("HX-Redirect", order.Link2)
	return c.String(http.StatusOK, "Order Created")
}

func OrderStatus(c echo.Context) error {
	//create fetch jorder fro gorm database
	var order models.Order
	var cust models.Customer

	orderID := c.Param("id")

	// Fetch the order from the GORM database
	if err := db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		// Handle the case where the order is not found
		return c.String(http.StatusNotFound, "Order not found")
	}

	sess, _ := session.Get("session", c)
	userID := sess.Values["user_id"]

	err := db.DB.Where("id = ?", userID).First(&cust).Error
	if err != nil {
		log.Println("func: OrderStatus - ", err)
	}

	// Update the PackageID
	db.DB.Model(&cust).Update("PackageID", order.PackageID)

	// Reload the Package association
	db.DB.Model(&cust).Association("Package").Find(&cust.Package)

	return c.Render(200, "order_status", order)
}

func ConfirmOrder(c echo.Context) error {
	//create fetch jorder fro gorm database
	var order models.Order
	var cust models.Customer

	orderID := c.Param("id")

	// Fetch the order from the GORM database
	log.Println("<< searching order ", orderID)
	if err := db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		// Handle the case where the order is not found
		log.Println("<< Order not found >>")
		return c.String(http.StatusNotFound, "Order not found")
	}

	orderStatus := paypal.CheckOrderStatus(order.PayId)

	sess, _ := session.Get("session", c)
	userID := sess.Values["user_id"]

	err := db.DB.Where("id = ?", userID).First(&cust).Error
	if err != nil {
		log.Println("func: OrderStatus - ", err)
	}

	// Update the PackageID
	db.DB.Model(&cust).Update("PackageID", order.PackageID)

	// Reload the Package association
	db.DB.Model(&cust).Association("Package").Find(&cust.Package)

	if orderStatus.Status == "APPROVED" {
		order.Status = true
		order.StatusMsg = orderStatus.Status
		db.DB.Save(&order)
	}

	user := map[string]interface{}{
		"name":  cust.Fname,
		"email": cust.Email,
		"subs":  cust.PackageID,
	}

	respData := map[string]interface{}{
		"user":  user,
		"order": order,
	}

	return c.Render(200, "order_status", respData)
}

func CancelOrder(c echo.Context) error {
	orderID := c.Param("id")

	//create fetch jorder fro gorm database
	var order models.Order

	// Fetch the order from the GORM database
	if err := db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		// Handle the case where the order is not found
		return c.String(http.StatusNotFound, "Order not found")
	}
	orderStatus := paypal.CheckOrderStatus(order.PayId)

	if orderStatus.Status == "APPROVED" {
		order.Status = true
		order.StatusMsg = orderStatus.Status
		db.DB.Save(order)
	}

	return c.Render(200, "order_status", order)
}
