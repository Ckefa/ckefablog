package handlers

import "github.com/labstack/echo/v4"

func HandleIndex(c echo.Context) error {
	return c.Render(200, "index", nil)
}

func ServerSideRendering(c echo.Context) error {
	return c.Render(200, "introtossr", nil)
}

func MentalHealth(c echo.Context) error {
	return c.Render(200, "mental_health", nil)
}
