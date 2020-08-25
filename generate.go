// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
)

const buildTags = "release"

func main() {
	err := vfsgen.Generate(assets, vfsgen.Options{
		BuildTags:    buildTags,
		VariableName: "assets",
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = vfsgen.Generate(templates, vfsgen.Options{
		BuildTags:    buildTags,
		VariableName: "templates",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
