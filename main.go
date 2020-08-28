package main

//go:generate rice embed-go

import (
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	assetsBox    = rice.MustFindBox("assets")
	templatesBox = rice.MustFindBox("templates")
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Configuration: %#v\n", config)

	t, err := loadTemplates(templatesBox)
	if err != nil {
		log.Fatalf("failed to load templates: %v", err)
	}
	log.Printf("Loaded templates%v\n", t.DefinedTemplates())

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Validator = NewValidator()
	e.Renderer = NewRenderer(t)

	generator, err := GeneratorFromFile(config.File, "-")
	if err != nil {
		log.Fatalf("failed to create generator: %v", err)
	}

	e.GET("/", defaultHandler(generator))
	e.GET("/:query", queryHandler(generator))
	e.GET("/_assets/*", echo.WrapHandler(http.StripPrefix("/_assets/", http.FileServer(assetsBox.HTTPBox()))))
	e.GET("/favicon.ico", redirectHandler("/_assets/favicons/favicon.ico"))

	e.Logger.Fatal(e.Start(config.Listen))
}
