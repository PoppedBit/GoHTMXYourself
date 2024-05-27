package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

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

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

func main() {

	// ENV
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port := os.Getenv("PORT")

	// Server
	e := echo.New()

	// Renderer
	e.Renderer = NewTemplates()

	// ???
	e.Use(middleware.Logger())

	count := Count{Count: 0}
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index", count)
	})

	e.POST("/count", func(c echo.Context) error {
		count.Count++
		return c.Render(http.StatusOK, "count", count)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
