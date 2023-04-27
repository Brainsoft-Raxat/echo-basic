package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}

	// Create a new Echo instance
	e := echo.New()

	// Use middleware to log all requests and recover from panics
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Define a route for the root path
	e.GET("/", func(c echo.Context) error {
		message := c.QueryParam("message")
		return c.String(http.StatusOK, fmt.Sprintf("Echo: %s", message))
	})

	// Start the server and listen on the specified port
	e.Logger.Fatal(e.Start(":" + port))
}
