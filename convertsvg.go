package convertsvg

import (
	"fmt"
	"image/png"
	"os"
	"path"
	"regexp"

	"bytes"
	"image"
	"io"

	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers"
)

type Renderer int

const (
	PNG Renderer = iota
	JPG
)

var rendererMap = map[Renderer]canvas.Writer{
	PNG: renderers.PNG(),
	JPG: renderers.JPEG(),
}

var extensionMap = map[Renderer]string{
	PNG: ".png",
	JPG: ".jpg",
}

func ConvertSvgFilesRecursive(sourcePath, destPath string, renderer Renderer) error {
	err := os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return err
	}
	dirEntries, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("error stat")
	}
	for _, dirEntry := range dirEntries {
		updatedPath := path.Join(sourcePath, dirEntry.Name())
		stat, err := os.Stat(updatedPath)
		if err != nil {
			return fmt.Errorf("error getting file status for %s: %w", updatedPath, err)
		}
		if stat.IsDir() {
			return ConvertSvgFilesRecursive(updatedPath, path.Join(destPath, dirEntry.Name()), renderer)
		}
		baseName := dirEntry.Name()
		ext := path.Ext(baseName)
		newExtension := extensionMap[renderer]
		replaceExtPattern, _ := regexp.Compile(regexp.QuoteMeta(ext) + "$")
		destBase := replaceExtPattern.ReplaceAllString(baseName, newExtension)
		updatedDestPath := path.Join(destPath, destBase)
		err = ConvertSvgFile(updatedPath, updatedDestPath, renderer)
		if err != nil {
			return err
		}
	}
	return nil
}

func ConvertSvgFile(sourcePath, destPath string, renderer Renderer) error {
	f, err := os.Open(sourcePath)

	if err != nil {
		return fmt.Errorf("error reading src file %s: %w", sourcePath, err)
	}
	img, err := ConvertSvg(f, renderer)
	if err != nil {
		return fmt.Errorf("error converting src file %s: %w", sourcePath, err)
	}
	newFile, err := os.OpenFile(destPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating destination file %s: %w", destPath, err)
	}
	err = png.Encode(newFile, img)
	if err != nil {
		return fmt.Errorf("error encoding png file %s: %w", destPath, err)
	}
	return nil
}

func ConvertSvg(reader io.Reader, renderer Renderer) (image.Image, error) {
	rendererFunc, ok := rendererMap[renderer]
	if !ok {
		return nil, fmt.Errorf("unrecognized renderer value %+v", renderer)
	}
	var err error
	c, err := canvas.ParseSVG(reader)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	err = c.Write(buf, rendererFunc)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(buf)
	if err != nil {
		return nil, err
	}
	return img, nil
}
