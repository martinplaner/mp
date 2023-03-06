package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	// FallbackQuery is used if no other query was specified
	FallbackQuery = "MP"
)

// RequestParams contains all request parameters
type RequestParams struct {
	Query string `param:"query" validate:"required,alphaunicode,min=1,max=20"`
}

// QueryResult contains the query as well as the result value of a single request
type QueryResult struct {
	Query  string `json:"query"`
	Result string `json:"result"`
}

// String implements fmt.Stringer used for plain text responses
func (r QueryResult) String() string {
	return r.Result
}

// Error holds useful information in case of an error
type Error struct {
	Message string `json:"message"`
	Cause   string `json:"cause"`
}

// Handler is responsible for handling requests using handler functions.
type Handler struct {
	Generator    Generator
	DefaultQuery string
}

func (h *Handler) rootHandler(c echo.Context) error {
	query := h.DefaultQuery

	if strings.TrimSpace(query) == "" {
		query = FallbackQuery
	}

	result, err := h.Generator.Generate(strings.ToUpper(query))
	if err != nil {
		return err
	}

	return Negotiate(http.StatusOK, "index.tmpl", &QueryResult{Query: query, Result: result}, c)
}

func (h *Handler) queryHandler(c echo.Context) error {
	req := new(RequestParams)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	query := strings.ToUpper(req.Query)
	result, err := h.Generator.Generate(query)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusNotFound,
			Message:  fmt.Sprintf("Could not generate word for query '%v'.", query),
			Internal: err,
		}
	}
	return Negotiate(http.StatusOK, "index.tmpl", &QueryResult{Query: query, Result: result}, c)
}

func redirectHandler(target string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, target)
	}
}
