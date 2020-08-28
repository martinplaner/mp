package main

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	// DefaultQuery is used if no other query was specified
	DefaultQuery = "MP"
)

// RequestParams contains all request parameters
type RequestParams struct {
	Query string `param:"query" validate:"required,alphaunicode,min=1,max=20"`
}

// QueryResult contains the query as well as the result value of a single request
type QueryResult struct {
	Query  string `json:"query"`
	Result string `json:"result" plain:"-"`
}

// Error holds useful information in case of an error
type Error struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

func defaultHandler(g *Generator) echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := g.Generate(strings.ToUpper(DefaultQuery))
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "index.tmpl", &QueryResult{Query: DefaultQuery, Result: result})
	}
}

func queryHandler(g *Generator) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(RequestParams)
		if err := c.Bind(req); err != nil {
			return err
		}
		if err := c.Validate(req); err != nil {
			return err
		}

		query := strings.ToUpper(req.Query)
		result, err := g.Generate(query)
		if err != nil {
			return err
		}

		return c.Render(http.StatusOK, "index.tmpl", &QueryResult{Query: query, Result: result})
	}
}

func redirectHandler(target string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, target)
	}
}
