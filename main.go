package main

//go:generate rice embed-go

import (
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Configuration: %#v\n", config)

	assetsBox, err := rice.FindBox("assets")
	if err != nil {
		log.Fatalf("failed to load assets box: %v", err)
	}
	log.Printf("Assets: appended=%v, embedded=%v", assetsBox.IsAppended(), assetsBox.IsEmbedded())

	templatesBox, err := rice.FindBox("templates")
	if err != nil {
		log.Fatalf("failed to templates assets box: %v", err)
	}
	log.Printf("Templates: appended=%v, embedded=%v", templatesBox.IsAppended(), templatesBox.IsEmbedded())

	t, err := loadTemplates(templatesBox)
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

	e.Validator = NewValidator()
	e.Renderer = NewRenderer(t)

	e.GET("/", defaultHandler(generator))
	e.GET("/:query", queryHandler(generator))
	e.GET("/_assets/*", echo.WrapHandler(http.StripPrefix("/_assets/", http.FileServer(assetsBox.HTTPBox()))))
	e.GET("/favicon.ico", redirectHandler("/_assets/favicons/favicon.ico"))

	e.Logger.Fatal(e.Start(config.Listen))
}
