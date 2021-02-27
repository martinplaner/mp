package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// MIME types
const (
	MIMEApplicationJSON = "application/json"
	MIMETextHTML        = "text/html"
	MIMETextPlain       = "text/plain"
	MIMEAny             = "*/*"
)

var (
	// ErrNotAcceptable is returned when no acceptable representation can be found.
	ErrNotAcceptable = &echo.HTTPError{
		Code:    http.StatusNotAcceptable,
		Message: "The target resource does not have a current representation that would be acceptable to the user agent.",
	}
)

type Validator struct {
	validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if ve, ok := err.(validator.ValidationErrors); ok {
		var sb strings.Builder
		for _, fe := range ve {
			sb.WriteString(fe.Field())
			sb.WriteString(":")
			sb.WriteString(fe.Tag())
		}
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  fmt.Sprintf("Request validation failed. [%v]", sb.String()),
			Internal: err,
		}
	}
	return err
}

func NewValidator() *Validator {
	return &Validator{validator: validator.New()}
}

// Renderer is a custom renderer implementation for template based HTML responses
type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.templates.ExecuteTemplate(w, name, data)
}

// Negotiate performs content negotiation using the "Accept" header and sends a response in the appropriate content type.
// If no content type can be negotiated, an HTTP 406 Not Acceptable error is returned.
func Negotiate(status int, name string, data interface{}, c echo.Context) error {
	accept := c.Request().Header.Get(http.CanonicalHeaderKey(echo.HeaderAccept))
	var accepts []string
	for _, a := range strings.Split(accept, ",") {
		accepts = append(accepts, strings.TrimSpace(a))
	}

	for _, a := range accepts {
		// Plain text (default)
		if strings.HasPrefix(a, MIMETextPlain) || strings.HasPrefix(a, MIMEAny) {
			if stringer, ok := data.(fmt.Stringer); ok {
				return c.String(status, stringer.String())
			}
		}

		// JSON
		if strings.HasPrefix(a, MIMEApplicationJSON) {
			return c.JSON(status, data)
		}

		// HTML
		if strings.HasPrefix(a, MIMETextHTML) {
			return c.Render(status, name, data)
		}
	}

	return ErrNotAcceptable
}

func NewRenderer(t *template.Template) *Renderer {
	return &Renderer{templates: t}
}
