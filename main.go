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
	"github.com/green-ecolution/tbz-csv-import-plugin/internal/importer/storage"
	"github.com/green-ecolution/tbz-csv-import-plugin/internal/server"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"
)

var version = "develop"

//go:embed all:ui/dist
var f embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	hostPathEnv := os.Getenv("HOST_PATH")

	pluginPath, err := url.Parse("http://localhost:8123/")
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

	db := sqlx.MustConnect("sqlite3", "file:import.db?cache=shared")
	importRepo := storage.NewImportRepositoryDB(db)

	if err = importRepo.Setup(); err != nil {
		slog.Error("Failed to migrate database", "error", err)
		panic(err)
	}

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
			URL:         fmt.Sprintf("%s/api", hostPathEnv),
			Description: "Green Ecolution API",
		},
	}
	clientCfg.Debug = true
	clientCfg.HTTPClient = oauthClient

	repo := storage.NewGreenEcolutionRepo(clientCfg)

	auth := context.WithValue(ctx, client.ContextOAuth2, oauthToken)
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
