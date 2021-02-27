package main

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed templates
var templates embed.FS

func loadTemplates() (*template.Template, error) {
	templateFS, err := fs.Sub(templates, "templates")
	if err != nil {
		return nil, err
	}

	t, err := template.ParseFS(templateFS, "*.tmpl")
	if err != nil {
		return nil, err
	}

	return t, nil
}
