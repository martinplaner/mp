package main

import (
	"embed"
	"io/fs"
)

//go:embed assets
var assets embed.FS

func loadAssets() (fs.FS, error) {
	assetsFS, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, err
	}

	return assetsFS, nil
}
