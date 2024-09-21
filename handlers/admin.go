package handlers

import (
	"log"
	"net/http"

	"github.com/Ckefa/ckefablog/db"
	"github.com/Ckefa/ckefablog/models"
	"github.com/labstack/echo/v4"
)

func Admin(c echo.Context) error {
	var respData = map[string]interface{}{}
	var orders []models.Order
	var pending, approved int

	if res := db.DB.Preload("customer").Find(&orders); res.Error != nil {
		log.Println("<< func: Admin - ", res.Error)
	}

	for i, order := range orders {
		// Count based on the Status field
		if order.Status {
			approved++
		} else {
			pending++
		}

		if err := db.DB.Where("id = ?", orders[i].CustomerID).Find(&orders[i].Customer).Error; err != nil {
			log.Println("customer order not found")
		}
	}

	respData["total"] = len(orders)
	respData["orders"] = orders
	respData["pending"] = pending
	respData["approved"] = approved

	return c.Render(http.StatusOK, "admin", respData)
}
