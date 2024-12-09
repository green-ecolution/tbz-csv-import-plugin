package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/green-ecolution/green-ecolution-backend/client"
	"github.com/green-ecolution/green-ecolution-backend/plugin"
	"github.com/green-ecolution/tbz-csv-import-plugin/internal/importer"
	"github.com/green-ecolution/tbz-csv-import-plugin/internal/server"
	"github.com/joho/godotenv"
)

var version = "develop"

//go:embed ui/dist/**/*
var f embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	hostPathEnv := os.Getenv("HOST_PATH")

	pluginPath, err := url.Parse("http://localhost:8080/")
	if err != nil {
		panic(err)
	}

	p := plugin.Plugin{
		Slug:           "csv-import",
		Name:           "CSV Import",
		Version:        version,
		Description:    "A plugin to import CSV files of trees from the TBZ Flensburg into the Green Ecolution system.",
		PluginHostPath: pluginPath,
	}

	http := server.NewServer(
		server.WithPort(8080),
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

	hostPath, err := url.Parse(hostPathEnv)
	if err != nil {
		panic(err)
	}

	worker, err := plugin.NewPluginWorker(
		plugin.WithHost(hostPath),
		plugin.WithPlugin(p),
		plugin.WithHostAPIVersion("v1"),
	)
	if err != nil {
		panic(err)
	}

	token, err := worker.Register(ctx, clientID, clientSecret)
	if err != nil {
		panic(err)
	}

	clientCfg := client.NewConfiguration()
	clientCfg.Servers = client.ServerConfigurations{
		{
			URL:         fmt.Sprintf("%s/api", hostPathEnv),
			Description: "Green Ecolution API",
		},
	}
	clientCfg.Debug = true
	clientCfg.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	repo := importer.NewGreenEcolutionRepo(clientCfg)
	auth := context.WithValue(context.Background(), client.ContextOAuth2, token)
	info, err := repo.GetInfo(auth)
	if err != nil {
		slog.Error("Error while getting app info", "error", err)
	}
	slog.Info("App info", "info", info)

	go func() {
		defer wg.Done()
		if err := worker.RunHeartbeat(ctx); err != nil {
			slog.Error("Failed to send heartbeat", "error", err)
		}
	}()

	wg.Wait()
}
