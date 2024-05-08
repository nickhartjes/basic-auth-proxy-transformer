package proxy

import "os"

type Settings struct {
	Port string
}

func GetSettings() *Settings {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Settings{
		Port: port,
	}
}
