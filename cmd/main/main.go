package main

import (
	"io"
	"log"
	"os"
	"text/template"

	"github.com/Ckefa/ckefablog/db"
	"github.com/Ckefa/ckefablog/handlers"
	"github.com/Ckefa/ckefablog/paypal"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
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

	// Initializing application vitals
	paypal.InitPayment()
	err = db.Init()
	if err != nil || db.DB == nil {
		log.Fatal("DB not initialized")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("Clinton#1234"))))
	e.Static("/", "static")

	e.Renderer = newTemplate()

	e.GET("/", handlers.HandleHome)
	e.GET("/blog", handlers.HandleIndex)

	e.POST("/subscribe", handlers.Subscribe)
	e.POST("/login", handlers.Login)
	e.POST("/signup", handlers.Register)
	e.POST("/pay", handlers.RequestOrder)

	e.GET("/login", handlers.HandleLogin)
	e.GET("/signup", handlers.Signup)
	e.GET("/logout", handlers.Logout)

	e.GET("/checkout/:pid", handlers.Checkout)

	e.GET("/order/confirm/:id", handlers.ConfirmOrder)
	e.GET("/oder/cancel/:id", handlers.CancelOrder)

	e.GET("/tech/golang-server-side-rendering", handlers.ServerSideRendering)
	e.GET("/tech/googlefi", handlers.GoogleFi)
	e.GET("/tech/future-of-remote-work", handlers.FutureOfRemoteWork)
	e.GET("/tech/ai-workplace-ethics", handlers.AiWorkplaceEthics)

	e.GET("/about", handlers.About)
	e.GET("/privacy-policy", handlers.PrivacyPolicy)
	e.GET("/terms-of-service", handlers.TermsOfService)

	e.Logger.Fatal(e.Start(PORT))
}
