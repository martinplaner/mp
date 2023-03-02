package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds the global program configuration (parsed from command line flags and environment variables).
type Config struct {
	Debug  bool
	File   string
	Listen string
	Files  *InputFiles
}

type Mode int64

const (
	Compound  Mode = 0
	Adjective Mode = 1
)

func (m Mode) String() string {
	switch m {
	case Compound:
		return "c"
	case Adjective:
		return "a"
	}

	panic("unexpected mode")
}

var modeMap = map[string]Mode{
	"c": Compound,
	"a": Adjective,
}

type InputFile struct {
	Language string
	Mode     Mode
	Filename string
}

type InputFiles []*InputFile

func (i *InputFiles) String() string {
	var b strings.Builder
	for _, f := range *i {
		b.WriteString(f.Language)
		b.WriteString(":")
		b.WriteString(f.Mode.String())
		b.WriteString(":")
		b.WriteString(f.Filename)
		b.WriteString(";")
	}
	return b.String()
}

func (i *InputFiles) Set(value string) error {
	var langArg string
	var mode string
	var filename string

	_, err := fmt.Sscanf(value, "%s:%s:%s", &langArg, &mode, &filename)
	if err != nil {
		return fmt.Errorf("could not parse file argument: %w", err)
	}

	*i = append(*i, &InputFile{
		Language: langArg,
		Mode:     Compound,
		Filename: filename,
	})

	return nil
}

func loadConfig() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	c := &Config{
		Files: &InputFiles{},
	}

	flag.BoolVar(&c.Debug, "debug", false, "Enable verbose debug mode")
	// flag.StringVar(&c.File, "file", "de:c:words_de.txt", "Path to word list (one word per line)")
	flag.StringVar(&c.Listen, "listen", fmt.Sprintf(":%s", port), "TCP address for the server to listen on, in the form 'host:port'")
	flag.Var(c.Files, "file", "TODO write input file usage")
	flag.Parse()

	return c, nil
}
