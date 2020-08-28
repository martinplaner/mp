package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"reflect"
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

// Renderer is a custom renderer implementation for different content types
type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	accept := c.Request().Header.Get(http.CanonicalHeaderKey(echo.HeaderAccept))
	var accepts []string
	for _, a := range strings.Split(accept, ",") {
		accepts = append(accepts, strings.TrimSpace(a))
	}

	for _, a := range accepts {
		// Plain text (default)
		if strings.HasPrefix(a, MIMETextPlain) || strings.HasPrefix(a, MIMEAny) {
			s, err := plainText(data)
			if err != nil {
				return err
			}
			_, err = w.Write([]byte(s))
			return nil
		}

		// JSON
		if strings.HasPrefix(a, MIMEApplicationJSON) {
			b, err := json.Marshal(data)
			if err != nil {
				return err
			}
			_, err = w.Write(b)
			return err
		}

		// HTML
		if strings.HasPrefix(a, MIMETextHTML) {
			return r.templates.ExecuteTemplate(w, name, data)
		}
	}

	return ErrNotAcceptable
}

func plainText(data interface{}) (string, error) {
	val := reflect.ValueOf(data)
	typ := val.Type()

	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = val.Type()
	}

	if typ.Kind() != reflect.Struct {
		return "", fmt.Errorf("data not of type struct (got %v)", typ.Kind())
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if _, ok := f.Tag.Lookup("plain"); ok {
			return val.Field(i).String(), nil
		}
	}

	return "", ErrNotAcceptable
}

func NewRenderer(t *template.Template) *Renderer {
	return &Renderer{templates: t}
}
