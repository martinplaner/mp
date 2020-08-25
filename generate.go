// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
)

func main() {
	err := vfsgen.Generate(assets, vfsgen.Options{
		BuildTags:    "release",
		VariableName: "assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
