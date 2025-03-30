package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/green-ecolution/green-ecolution-backend/pkg/client"
	"github.com/green-ecolution/green-ecolution-backend/pkg/plugin"
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
		log.Println("Error loading .env file")
	}

	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	p := plugin.Plugin{
		Slug:           "csv-import",
		Name:           "CSV Import",
		Version:        version,
		Description:    "A plugin to import CSV files of trees from the TBZ Flensburg into the Green Ecolution system.",
		PluginHostPath: cfg.PluginPath,
	}

	worker, err := plugin.NewPluginWorker(
		plugin.WithHost(cfg.HostPath),
		plugin.WithPlugin(p),
		plugin.WithHostAPIVersion("v1"),
		plugin.WithClientID(cfg.ClientID),
		plugin.WithClientSecret(cfg.ClientSecret),
	)
	if err != nil {
		panic(err)
	}

	oauthClient := authClient(ctx, worker)
	clientCfg := client.NewConfiguration()
	clientCfg.Servers = client.ServerConfigurations{
		{
			URL:         fmt.Sprintf("%s/api", cfg.HostPath),
			Description: "Green Ecolution API",
		},
	}
	clientCfg.Debug = true
	clientCfg.HTTPClient = oauthClient

	geClient := NewGreenEcolutionRepo(clientCfg, p.Slug)

	fSub, err := fs.Sub(f, "ui/dist")
	if err != nil {
		panic(err)
	}

	var serverPort int
	if cfg.PluginPath.Port() != "" {
		serverPort, err = strconv.Atoi(cfg.PluginPath.Port())
		if err != nil {
			panic(err)
		}
	}

	server := NewServer(
		WithPort(serverPort),
		WithPluginFS(fSub),
		WithPlugin(p),
		WithVersion(version),
		WithClient(geClient),
	)

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		if err = server.Run(ctx); err != nil {
			slog.Error("error while running http server", "error", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := worker.RunHeartbeat(ctx); err != nil {
			slog.Error("failed to send heartbeat", "error", err)
			cleanup(worker)
		}
	}()

	go func() {
		defer wg.Done()
		<-ctx.Done()
		cleanup(worker)
	}()

	wg.Wait()
}

func authClient(ctx context.Context, worker *plugin.PluginWorker) *http.Client {
	token, err := worker.Register(ctx)
	if err != nil {
		panic(err)
	}

	oauthToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
		TokenType:    token.TokenType,
	}

	return oauth2.NewClient(ctx, NewTokenSource(worker.RefreshToken, oauthToken))
}

func cleanup(worker *plugin.PluginWorker) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := worker.Unregister(timeoutCtx); err != nil {
		slog.Error("failed to unregister plugin", "error", err)
	}

	os.Exit(1)
}
