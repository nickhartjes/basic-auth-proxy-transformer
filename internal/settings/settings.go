package settings

import (
	"log"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type RedisSettings struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
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

type Settings struct {
	Port  string        `mapstructure:"port"`
	Cache CacheSettings `mapstructure:"cache"`
}

func loadDefaultSettings() {
	viper.SetDefault("port", "8080")
	viper.SetDefault("cache.enabled", true)
	viper.SetDefault("cache.cache_type", "ristretto")
	viper.SetDefault("cache.ristretto.num_counters", 1000)
	viper.SetDefault("cache.ristretto.max_cost", 100)
	viper.SetDefault("cache.ristretto.buffer_items", 64)
	viper.SetDefault("cache.redis.addr", "localhost:6379")
	viper.SetDefault("cache.redis.password", "")
	viper.SetDefault("cache.redis.db", 0)
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

	// Initialize logging based on an environment variable or other criteria
	if debug := os.Getenv("DEBUG"); debug == "true" {
		log.SetFlags(log.Ldate | log.Lmicroseconds)
	} else {
		log.SetFlags(log.Ldate | log.LstdFlags)
	}

	return &settings
}
