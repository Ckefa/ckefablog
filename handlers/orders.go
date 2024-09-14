package handlers

import (
	"log"
	"net/http"

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
