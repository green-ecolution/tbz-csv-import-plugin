package server

import (
	"context"
	"embed"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/green-ecolution/green-ecolution-backend/pkg/plugin"
)

type ServerConfig struct {
	port     int
	plugin   plugin.Plugin
	pluginFS embed.FS
}

type Server struct {
	cfg *ServerConfig
}

type ServerOption func(*ServerConfig)

func WithPort(port int) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.port = port
	}
}

func WithPluginFS(pluginFS embed.FS) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.pluginFS = pluginFS
	}
}

func WithPlugin(plugin plugin.Plugin) ServerOption {
	return func(cfg *ServerConfig) {
		cfg.plugin = plugin
	}
}

var defaultServerConfig = &ServerConfig{
	port: 8080,
}

func NewServer(opts ...ServerOption) *Server {
	cfg := defaultServerConfig
	for _, opt := range opts {
		opt(cfg)
	}
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Run(ctx context.Context) error {
	app := fiber.New(fiber.Config{})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World! This is the plugin server for " + s.cfg.plugin.Name)
	})
	app.Mount("/", servePlugin(s.cfg.pluginFS))

	go func() {
		<-ctx.Done()
		fmt.Println("Shutting down HTTP Server")
		if err := app.Shutdown(); err != nil {
			fmt.Println("Error while shutting down HTTP Server:", err)
		}
	}()

	return app.Listen(fmt.Sprintf(":%d", s.cfg.port))
}
