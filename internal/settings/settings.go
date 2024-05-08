package settings

import (
	"log"
	"log/slog"
	"os"
)

type Settings struct {
	Port string
}

func GetSettings() *Settings {
	debug := os.Getenv("DEBUG")
	if debug == "" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		log.SetFlags(log.Ldate | log.Lmicroseconds)
	} else {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		log.SetFlags(log.Ldate | log.Lmicroseconds)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Settings{
		Port: port,
	}
}
