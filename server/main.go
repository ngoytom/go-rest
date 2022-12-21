package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// stored in memory
type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
	Body  string `json:"body"`
}

func main() {
	app := fiber.New()

	//fix CORS issue
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// create slice of Todo
	todos := []Todo{}

	/* Test server
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	*/

	// Create new Todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		// Checks JSON data in Todo and checks if there are any error
		if err := c.BodyParser(todo); err != nil {
			return err
		}

		// Give unique ID
		todo.ID = len(todos) + 1

		todos = append(todos, *todo)

		// Returns our slice of todos
		return c.JSON(todos)
	})

	// Update Partially so PATCH request rather than PUT
	app.Patch("/api/todos/:id/done", func(c *fiber.Ctx) error {
		// Get ID out of url and conver to integer
		id, err := c.ParamsInt("id")

		// If error is found
		if err != nil {
			return c.Status(401).SendString("Invalid id")
		}

		// Loop through all of todos and if id match, mark status as true
		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = true
				break
			}
		}

		return c.JSON(todos)
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		// return todos in JSON format
		return c.JSON(todos)
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		// get id from url as int
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(401).SendString("Invalid ID")
		}

		// delete todo by id
		for i, t := range todos {
			if t.ID == id {
				todos = append(todos[:i], todos[i+1:]...)
				break
			}
		}
		return c.JSON(todos)
	})

	log.Fatal(app.Listen(":4000"))
}
