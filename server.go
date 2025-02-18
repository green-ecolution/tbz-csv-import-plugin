package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/green-ecolution/green-ecolution-backend/pkg/plugin"
)

type ServerConfig struct {
	port     int
	client   *GreenEcolutionClient
	plugin   plugin.Plugin
	pluginFS fs.FS
	version  string
}

type Server struct {
	cfg *ServerConfig
}

type ServerOption func(*ServerConfig)

func WithPort(port int) ServerOption {
	return func(sc *ServerConfig) { sc.port = port }
}

func WithClient(client *GreenEcolutionClient) ServerOption {
	return func(sc *ServerConfig) { sc.client = client }
}

func WithPlugin(plugin plugin.Plugin) ServerOption {
	return func(sc *ServerConfig) { sc.plugin = plugin }
}

func WithPluginFS(fs fs.FS) ServerOption {
	return func(sc *ServerConfig) { sc.pluginFS = fs }
}

func WithVersion(version string) ServerOption {
	return func(sc *ServerConfig) { sc.version = version }
}

var defaultServerConfig = &ServerConfig{
	port:    8080,
	version: "develop",
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
	r := chi.NewRouter()

	r.Get("/", s.handleHelloWorld)
	r.Get("/*", s.handleFileSystem)

	r.Post("/upload", s.handleCsvUpload)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", s.cfg.port),
		Handler: r,
	}

	go func() {
		<-ctx.Done()
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := server.Shutdown(timeoutCtx); err != nil {
			slog.Error("failed to shutdown http server", "error", err)
		}
	}()

	return server.ListenAndServe()
}

func (s *Server) handleHelloWorld(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello World"))
}

func (s *Server) handleFileSystem(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
	fs := http.StripPrefix(pathPrefix, http.FileServerFS(s.cfg.pluginFS))
	fs.ServeHTTP(w, r)
}

func (s *Server) handleCsvUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	r.ParseForm()

	csvFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("failed to read file: error: %s", err)))
		return
	}

	csvConverter := NewCSVConverter(fileHeader, csvFile)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("unsupported file: error: %s", err)))
		return
	}

	convertedTrees, err := csvConverter.Convert(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to convert csv file: error: %s", err)))
		return
	}

	syncTrees := NewSyncTrees(convertedTrees, s.cfg.client)
	importedTrees, err := syncTrees.Sync(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("failed to sync trees to green ecolution backend: error: %s", err)))
		return
	}

	encode := json.NewEncoder(w)
	encode.Encode(importedTrees)
}
