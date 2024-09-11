package handlers

import (
	"net/http"

	"github.com/Ckefa/ckefablog.git/db"
	"github.com/Ckefa/ckefablog.git/models"
	"github.com/Ckefa/ckefablog.git/paypal"
	"github.com/labstack/echo/v4"
)

func Checkout(c echo.Context) error {
	return c.Render(200, "checkout", nil)
}

func RequestOrder(c echo.Context) error {
	amt := c.FormValue("amount")
	order := models.NewOrder(amt)
	resp := paypal.CreateOrder(order)
	db.DB.Save(order)

	return c.Redirect(http.StatusTemporaryRedirect, resp.Links[1].Href)
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

	return c.Render(200, "order_status", map[string]interface{}{
		"order": order,
	})
}

func ConfirmOrder(c echo.Context) error {
	orderID := c.Param("id")
	status := paypal.CheckOrderStatus(orderID)

	//create fetch jorder fro gorm database
	var order models.Order

	// Fetch the order from the GORM database
	if err := db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		// Handle the case where the order is not found
		return c.String(http.StatusNotFound, "Order not found")
	}

	if status.Status == "APPROVED" {
		order.Status = true
		db.DB.Save(order)
	}
	return c.Render(200, "order_status", map[string]interface{}{
		"order": order,
	})
}

func CancelOrder(c echo.Context) error {
	orderID := c.Param("id")
	status := paypal.CheckOrderStatus(orderID)

	//create fetch jorder fro gorm database
	var order models.Order

	// Fetch the order from the GORM database
	if err := db.DB.Where("id = ?", orderID).First(&order).Error; err != nil {
		// Handle the case where the order is not found
		return c.String(http.StatusNotFound, "Order not found")
	}

	if status.Status == "APPROVED" {
		order.Status = true
		db.DB.Save(order)
	}

	return c.Render(200, "order_status", map[string]interface{}{
		"order": order,
	})
}
