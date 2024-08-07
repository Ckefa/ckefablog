package handlers

import "github.com/labstack/echo/v4"

func ServerSideRendering(c echo.Context) error {
	return c.Render(200, "introtossr", "")
}
