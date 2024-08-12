package main

import (
	"github.com/gofiber/fiber/v2"

	"html/template"
)

func main() {

	app := fiber.New()
	// Create a new template engine instance
	engine := template.New("views")
	// Parse your template files
	templateFile := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>{{.Title}}</title>
    </head>
    <body>
        <h1>{{.Content}}</h1>
    </body>
    </html>
    `
	t, err := engine.Parse(templateFile)
	if err != nil {
		// Handle the error
		return
	}

	// Define a route to render the template
	app.Get("/", func(c *fiber.Ctx) error {
		data := struct {
			Title   string
			Content string
		}{
			Title:   "My Page",
			Content: "Welcome to my website!",
		}

		// Add a condition to display additional content
		if true {
			data.Content = "You're a registered user!"
		}
		return t.Execute(c.Response().BodyWriter(), data)
	})
	app.Listen(":3000")

}
