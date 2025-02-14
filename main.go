package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/green-ecolution/green-ecolution-backend/pkg/client"
	"github.com/green-ecolution/green-ecolution-backend/pkg/plugin"
	"github.com/green-ecolution/tbz-csv-import-plugin/internal/server"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	version = "develop"
	cfg     Config
)

//go:embed all:ui/dist
var f embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	p := plugin.Plugin{
		Slug:           "csv-import",
		Name:           "CSV Import",
		Version:        version,
		Description:    "A plugin to import CSV files of trees from the TBZ Flensburg into the Green Ecolution system.",
		PluginHostPath: cfg.PluginPath,
	}

	http := server.NewServer(
		server.WithPort(8123),
		server.WithPluginFS(f),
		server.WithPlugin(p),
		server.WithVersion(version),
	)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	//wg.Add(2)
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err = http.Run(ctx); err != nil {
			slog.Error("Error while running http server", "error", err)
		}
	}()

	worker, err := plugin.NewPluginWorker(
		plugin.WithHost(cfg.HostPath),
		plugin.WithPlugin(p),
		plugin.WithHostAPIVersion("v1"),
	)
	if err != nil {
		panic(err)
	}

	token, err := worker.Register(ctx, cfg.ClientID, cfg.ClientSecret)
	if err != nil {
		panic(err)
	}

	oauthToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		ExpiresIn:    token.ExpiresIn,
		TokenType:    "Bearer",
	}
	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(oauthToken))
	clientCfg := client.NewConfiguration()
	clientCfg.Servers = client.ServerConfigurations{
		{
			URL:         fmt.Sprintf("%s/api", cfg.HostPath),
			Description: "Green Ecolution API",
		},
	}
	clientCfg.Debug = true
	clientCfg.HTTPClient = oauthClient

	go func() {
		defer wg.Done()
		if err := worker.RunHeartbeat(ctx); err != nil {
			slog.Error("Failed to send heartbeat", "error", err)
		}
	}()

	wg.Wait()
}
