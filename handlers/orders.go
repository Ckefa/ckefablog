package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Ckefa/ckefablog/db"
	"github.com/Ckefa/ckefablog/models"
	"github.com/Ckefa/ckefablog/paypal"
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
	amt := c.FormValue("amount")
	log.Println("Paying $", amt)
	amt_float, err := strconv.ParseFloat(amt, 64)
	pid := models.GetPid(amt_float)

	order := models.NewOrder(amt, pid)
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
	orderID := c.Param("id")
	var order models.Order

	// Fetch the order from the GORM database
	if err := db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		// Handle the case where the order is not found
		return c.String(http.StatusNotFound, "Order not found")
	}

	return c.Render(200, "order_status", order)
}

func ConfirmOrder(c echo.Context) error {
	orderID := c.Param("id")

	//create fetch jorder fro gorm database
	var order models.Order

	// Fetch the order from the GORM database
	log.Println("<< searching order ", orderID)
	if err := db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		// Handle the case where the order is not found
		log.Println("<< Order not found >>")
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
