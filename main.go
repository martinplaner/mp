package main

//go:generate go run generate.go assets.go

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type CustomRenderer struct {
	templates *template.Template
}

func (r *CustomRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	accept := c.Request().Header.Get("Accept")

	if accept == "text/plain" {
		return c.String(http.StatusOK, "plain")
	}

	return c.JSON(http.StatusOK, data)
}

func main() {

	config, err := loadConfig()
	if err != nil {
		log.Fatal("failed to load config", zap.Error(err))
	}

	log.Printf("Configuration: %#v\n", config)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Validator = &CustomValidator{validator: validator.New()}
	e.Renderer = &CustomRenderer{}

	e.Static("/_assets", "assets")
	e.GET("/:query", queryHandler)

	// generator, err := GeneratorFromFile(config.File, "-")
	// if err != nil {
	// 	logger.Fatal("failed to create generator", zap.Error(err))
	// }

	e.Logger.Fatal(e.Start(config.Listen))
}

type Request struct {
	Query string `param:"query" validate:"required,alphaunicode,min=1,max=20"`
}

func queryHandler(c echo.Context) error {
	req := new(Request)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	if strings.Contains(req.Query, "a") {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "query must not contain 'a'",
		}
	}

	return c.Render(http.StatusOK, "index", req.Query)
	// return c.JSON(http.StatusOK, req.Query)
}

func loadTemplate(name string) (*template.Template, error) {
	file, err := assets.Open("templates/" + name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	h, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	t, err := template.New(name).Parse(string(h))
	if err != nil {
		return nil, err
	}

	return t, nil
}
