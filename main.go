package main

//go:generate rice embed-go

import (
	"log"
	"net/http"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Configuration: %#v\n", config)

	assetsFS, err := loadAssets()
	if err != nil {
		log.Fatalf("failed to load assets: %v", err)
	}

	t, err := loadTemplates()
	if err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}
	log.Printf("Loaded templates%v\n", t.DefinedTemplates())

	generator, err := GeneratorFromFile(config.File, "-")
	if err != nil {
		log.Fatalf("failed to create generator: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	p := prometheus.NewPrometheus("mp", nil)
	p.MetricsPath = "_metrics"
	p.Use(e)

	e.Validator = NewValidator()
	e.Renderer = NewRenderer(t)

	e.GET("/", defaultHandler(generator))
	e.GET("/:query", queryHandler(generator))
	e.GET("/_assets/*", echo.WrapHandler(http.StripPrefix("/_assets/", http.FileServer(http.FS(assetsFS)))))
	e.GET("/favicon.ico", redirectHandler("/_assets/favicons/favicon.ico"))

	e.Logger.Fatal(e.Start(config.Listen))
}
