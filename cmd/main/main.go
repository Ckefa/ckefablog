package main

import (
	"io"
	"log"
	"os"
	"text/template"

	"github.com/Ckefa/ckefablog.git/db"
	"github.com/Ckefa/ckefablog.git/handlers"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("contents/**/*.html")),
	}
}

func main() {
	println("Starting app .....")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load environment variables")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("<<Port not Configured>>")
	}

	err = db.Init()
	if err != nil || db.DB == nil {
		log.Fatal("DB not initialized")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Static("/", "static")

	e.Renderer = newTemplate()

	e.GET("/", handlers.HandleIndex)
	e.POST("/subscribe", handlers.Signup)

	e.GET("/tech/golang-server-side-rendering", handlers.ServerSideRendering)
	e.GET("/tech/googlefi", handlers.GoogleFi)
	e.GET("/tech/future-of-remote-work", handlers.FutureOfRemoteWork)
	e.GET("/tech/ai-workplace-ethics", handlers.AiWorkplaceEthics)

	e.GET("/about", handlers.About)
	e.GET("/privacy-policy", handlers.PrivacyPolicy)
	e.GET("/terms-of-service", handlers.TermsOfService)

	e.Logger.Fatal(e.Start(PORT))
}
