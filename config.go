package main

import (
	"net/url"
)

type Config struct {
	PluginSlug   string   `env:"GE_PLUGIN_SLUG"`
	PluginPath   *url.URL `env:"GE_PLUGIN_PATH"`
	HostPath     *url.URL `env:"GE_HOST_PATH"`
	ClientID     string   `env:"GE_CLIENT_ID"`
	ClientSecret string   `env:"GE_CLIENT_SECRET"`
	CsvHeaders   []string `env:"GE_CSV_HEADERS" envSeparator:","`
	CsvUsedEpsg  string   `env:"GE_CSV_USED_EPSG"`
	CsvToEpsg    string   `env:"GE_CSV_To_EPSG"`
}
