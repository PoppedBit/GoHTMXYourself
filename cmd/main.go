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

type Contact struct {
	Name  string
	Email string
}

func newContact(name, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

type Contacts = []Contact

type ContactsData struct {
	Contacts Contacts
}

func (c *ContactsData) hasEmail(email string) bool {
	for _, contact := range c.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func newContactsDate() ContactsData {
	return ContactsData{
		Contacts: Contacts{
			newContact("John Doe", "john.doe@email.com"),
			newContact("Jane Doe", "jane.doe@email.com"),
		},
	}
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

	// Contacts
	contacts := newContactsDate()
	e.GET("/contacts", func(c echo.Context) error {
		return c.Render(http.StatusOK, "contacts", contacts)
	})
	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if contacts.hasEmail(email) {
			return c.String(http.StatusBadRequest, "Email already exists!")
		}

		contact := newContact(name, email)

		contacts.Contacts = append(contacts.Contacts, contact)

		return c.Render(http.StatusOK, "contacts-table-row", contact)
	})

	e.Logger.Fatal(e.Start(":" + port))
}
