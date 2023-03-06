package main

import (
	"flag"
	"fmt"
	"os"
)

// Config holds the global program configuration (parsed from command line flags and environment variables).
type Config struct {
	Debug        bool
	DefaultQuery string
	File         string
	Listen       string
	Mode         Mode
	Once         string
}

//go:generate stringer -type Mode -output config_string.go

// Mode defines the word generation mode.
// Either compound nouns (default) or a single noun prefixed by a list of adjectives.
type Mode int64

const (
	Compound  Mode = 0
	Adjective Mode = 1
)

func (m Mode) Set(value string) error {
	switch value {
	case "a", "adjective":
		m = Adjective
		return nil
	case "c", "compound":
		m = Compound
		return nil
	}

	return fmt.Errorf("invalid mode value: %v", value)
}

func loadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c := &Config{}

	flag.BoolVar(&c.Debug, "debug", false, "Enable verbose debug mode")
	flag.StringVar(&c.DefaultQuery, "default", "MP", "Default fallback query term, if not provided.")
	flag.StringVar(&c.File, "file", "words.txt", "Path to word list (one word per line, with optional word type)")
	flag.StringVar(&c.Listen, "listen", fmt.Sprintf(":%s", port), "TCP address for the server to listen on, in the form 'host:port'")
	flag.Var(&c.Mode, "mode", "Word generation mode, one of a|adjective|c|compound")
	flag.StringVar(&c.Once, "once", "", "Run generation once with the given query and print result, then quit")
	flag.Parse()

	return c, nil
}
