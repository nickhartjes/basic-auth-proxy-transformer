package internal

import (
	"github.com/spf13/viper"
	"log"
	"log/slog"
)

type RedisSettings struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	Database int    `mapstructure:"db"`
}

type RistrettoSettings struct {
	NumCounters int64 `mapstructure:"num_counters"`
	MaxCost     int64 `mapstructure:"max_cost"`
	BufferItems int64 `mapstructure:"buffer_items"`
}

type CacheSettings struct {
	Enabled   bool              `mapstructure:"enabled"`
	CacheType string            `mapstructure:"cache_type"`
	Ristretto RistrettoSettings `mapstructure:"ristretto"`
	Redis     RedisSettings     `mapstructure:"redis"`
}

type OAuth2Settings struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	TokenEndpoint string `mapstructure:"token_endpoint"`
	ClientID      string `mapstructure:"client_id"`
	ClientSecret  string `mapstructure:"client_secret"`
}

type Settings struct {
	Port             string         `mapstructure:"port"`
	TargetHeaderName string         `mapstructure:"target_header_name"`
	Debug            bool           `mapstructure:"debug"`
	Cache            CacheSettings  `mapstructure:"cache"`
	OAuth2           OAuth2Settings `mapstructure:"oauth2"`
}

func loadDefaultSettings() {
	viper.SetDefault("port", "8080")
	viper.SetDefault("debug", false)
	viper.SetDefault("target_header_name", "X-Target-URL")
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.cache_type", "ristretto")
	viper.SetDefault("cache.ristretto.num_counters", 1000)
	viper.SetDefault("cache.ristretto.max_cost", 100)
	viper.SetDefault("cache.ristretto.buffer_items", 64)
	viper.SetDefault("cache.redis.host", "http://localhost")
	viper.SetDefault("cache.redis.port", 6379)
	viper.SetDefault("cache.redis.password", "")
	viper.SetDefault("cache.redis.database", 0)
	viper.SetDefault("oauth2.host", "localhost")
	viper.SetDefault("oauth2.port", "8090")
	viper.SetDefault("oauth2.token_endpoint", "/realms/example/protocol/openid-connect/token")
	viper.SetDefault("oauth2.client_id", "my-client")
	viper.SetDefault("oauth2.client_secret", "my-client-secret")
}

func configureEnvironmentOverrides() {
	viper.AutomaticEnv() // Use environment variables to override
}

func configureConfigFile() {
	viper.SetConfigName("config") // Name of the config file (without extension)
	viper.SetConfigType("toml")   // Type of the config file
	viper.AddConfigPath(".")      // Look for the config file in the current directory

	err := viper.ReadInConfig() // Read the config file if it exists
	if err != nil {
		slog.Warn("Warning: unable to read config file, using defaults: %v", err)
	}
}

// GetSettings retrieves the settings from the configuration file and environment variables
func GetSettings() *Settings {
	// Load default settings, file settings, and environment overrides
	loadDefaultSettings()
	configureConfigFile()
	configureEnvironmentOverrides()

	// Retrieve settings into a struct
	var settings Settings
	err := viper.Unmarshal(&settings)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}
	return &settings
}
