package main

import (
	"flag"
	"fmt"
	"os"
)

// Config holds the global program configuration (parsed from command line flags and environment variables).
type Config struct {
	File   string
	Listen string
}

func loadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c := &Config{}

	flag.StringVar(&c.File, "file", "words.txt", "Path to word list (one word per line)")
	flag.StringVar(&c.Listen, "listen", fmt.Sprintf(":%s", port), "TCP address for the server to listen on, in the form 'host:port'")
	flag.Parse()

	return c, nil
}
