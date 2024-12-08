package server

import (
	"embed"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func servePlugin(f embed.FS) *fiber.App {
	app := fiber.New()

	app.Use(filesystem.New(filesystem.Config{
		Root:       http.FS(f),
		PathPrefix: "ui/dist",
		Browse:     true,
	}))

	return app
}
