package handlers

import "github.com/labstack/echo/v4"

func PrivacyPolicy(c echo.Context) error {
	return c.Render(200, "privacy_policy", nil)
}

func TermsOfService(c echo.Context) error {
	return c.Render(200, "terms_of_service", nil)
}

func About(c echo.Context) error {
	return c.Render(200, "about", nil)
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

/*
	Lifestyle
*/

func MentalHealth(c echo.Context) error {
	return c.Render(200, "mental_health", nil)
}

func FrugalLivingTips(c echo.Context) error {
	return c.Render(200, "frugal_living", nil)
}
