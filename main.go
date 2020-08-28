package main

//go:generate rice embed-go

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	assets = rice.MustFindBox("assets")
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	log.Printf("Configuration: %#v\n", config)

	t, err := loadTemplates()
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

	err = assets.Walk("/", func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})

	e.GET("/", defaultHandler(generator))
	e.GET("/:query", queryHandler(generator))
	e.GET("/_assets/*", echo.WrapHandler(http.StripPrefix("/_assets/", http.FileServer(assets.HTTPBox()))))
	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/_assets/favicons/favicon.ico")
	})

	e.Logger.Fatal(e.Start(config.Listen))
}

func loadTemplates() (*template.Template, error) {
	t := template.New("")
	err := assets.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		baseName := filepath.Base(path)
		if _, err := filepath.Match("*.tmpl", baseName); err != nil {
			return err
		}
		file, err := assets.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		h, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		t, err = template.New(baseName).Parse(string(h))
		if err != nil {
			return err
		}
		return nil
	})
	return t, err
}
