package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	if config.Debug {
		log.Printf("Configuration: %#v\n", config)
	}

	var generator Generator

	if config.Mode == Adjective {
		generator, err = AdjectiveGeneratorFromFile(config.File, " ")
	} else {
		generator, err = CompoundGeneratorFromFile(config.File, "-")
	}

	if err != nil {
		log.Fatalf("failed to create generator: %v", err)
	}

	if config.Once != "" {
		output, err := generator.Generate(strings.ToUpper(config.Once))
		if err != nil {
			log.Fatalf("failed to generate output: %v", err)
		}
		fmt.Println(output)
		return
	}

	assetsFS, err := loadAssets()
	if err != nil {
		log.Fatalf("failed to load assets: %v", err)
	}

	t, err := loadTemplates()
	if err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}
	log.Printf("Loaded templates%v\n", t.DefinedTemplates())

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	p := prometheus.NewPrometheus("mp", nil)
	p.MetricsPath = "_metrics"
	p.Use(e)

	e.Validator = NewValidator()
	e.Renderer = NewRenderer(t)

	handler := Handler{
		Generator:    generator,
		DefaultQuery: config.DefaultQuery,
	}

	e.GET("/", handler.rootHandler)
	e.GET("/:query", handler.queryHandler)
	e.GET("/_assets/*", echo.WrapHandler(http.StripPrefix("/_assets/", http.FileServer(http.FS(assetsFS)))))
	e.GET("/favicon.ico", redirectHandler("/_assets/favicons/favicon.ico"))

	e.Logger.Fatal(e.Start(config.Listen))
}
