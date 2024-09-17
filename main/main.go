package main

import (
	"convertsvg"
	"flag"
	"fmt"
)

var rendererFlagMap = map[string]convertsvg.Renderer{
	"png": convertsvg.PNG,
	"jpg": convertsvg.JPG,
}

func main() {
	recursive := flag.Bool("recursive", false, "walk directories recursively")
	rendererFlag := flag.String("renderer", "png", "renderer")
	sourcePath := flag.String("source", "", "source file")
	destPath := flag.String("dest", "", "source file")
	flag.Parse()

	renderer, ok := rendererFlagMap[*rendererFlag]
	if !ok {
		panic(fmt.Sprintf("unrecognized renderer %q", *rendererFlag))
	}
	if *recursive {
		err := convertsvg.ConvertSvgFilesRecursive(*sourcePath, *destPath, renderer)
		if err != nil {
			panic(err)
		}
		return
	}
	err := convertsvg.ConvertSvgFile(*sourcePath, *destPath, renderer)
	if err != nil {
		panic(err)
	}
}
