package main

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	rice "github.com/GeertJohan/go.rice"
)

func loadTemplates(box *rice.Box) (*template.Template, error) {
	t := template.New("")
	err := box.Walk("/", func(path string, info os.FileInfo, err error) error {
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
		file, err := box.Open(path)
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
