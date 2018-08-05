// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/wpdirectory/wpdir/internal/data"
)

func main() {
	err := vfsgen.Generate(data.Assets, vfsgen.Options{
		Filename:     "../../internal/data/vfsdata.go",
		PackageName:  "data",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
