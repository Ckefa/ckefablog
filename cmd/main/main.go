package main

import (
	"io"
	"text/template"

	"github.com/Ckefa/ckefablog.git/handlers"
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
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Unable to load environment variables")
	// }
	//
	// addr := os.Getenv("laddr")
	// if addr == "" {
	// 	log.Fatal("Environment variable 'laddr' not set")
	// }
	//
	// dsn := os.Getenv("dsn")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Static("/", "static")

	e.Renderer = newTemplate()

	e.GET("/", handlers.HandleIndex)
	e.GET("/golang-server-side-rendering", handlers.ServerSideRendering)
	e.GET("/prioritizing-mental-health", handlers.MentalHealth)

	e.Logger.Fatal(e.Start(":3000"))
}
