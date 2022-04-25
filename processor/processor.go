package processor

import (
	"bytes"
	"converter/pdf"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const (
	FormatJPEG    = ".jpg"
	FormatPNG     = ".png"
	FormatSVG     = ".svg"
	DefaultFormat = FormatJPEG
)

type Entity struct {
	SourceFile        string
	SourceFormat      string
	TargetFolder      string
	ConvertFileFormat string
	PrefixFileName    string
	TargetFile        string
	Converter         pdf.Converter
	bSourceFile       []byte
	bConvertFiles     [][]byte
}

func (e *Entity) Convert() error {

	if err := e.readSourceFile(); err != nil {
		return err
	}
	if err := e.convertSource(); err != nil {
		return err
	}
	if err := e.saveAllConvertsFiles(); err != nil {
		return err
	}
	return nil
}

func (e *Entity) readSourceFile() error {
	b, err := ioutil.ReadFile(e.SourceFile)
	if err == nil {
		e.bSourceFile = b
	}
	return err
}

func (e *Entity) convertSource() error {
	var err error
	var ib [][]byte
	switch e.ConvertFileFormat {
	case FormatJPEG:
		ib, err = e.Converter.Convert(e.bSourceFile, pdf.ImageFormatJPEG)
	case FormatPNG:
		ib, err = e.Converter.Convert(e.bSourceFile, pdf.ImageFormatPNG)
	default:
		ib, err = e.Converter.Convert(e.bSourceFile, pdf.ImageFormatJPEG)
	}
	if err == nil {
		e.bConvertFiles = ib
	}
	return err
}

func (e Entity) getConvertFileByArray(n int) (image.Image, string, error) {
	return image.Decode(bytes.NewReader(e.bConvertFiles[n]))
}

func (e *Entity) createFile(img image.Image, fileName string) error {
	f, err := os.Create(filepath.Join(fileName))
	if err != nil {
		return err
	}
	defer f.Close()
	err = jpeg.Encode(f, img, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		return err
	}
	return nil
}

func (e *Entity) saveAllConvertsFiles() error {
	if err := os.MkdirAll(e.TargetFolder, os.ModePerm); err != nil {
		return err
	}
	for n := 0; n < len(e.bConvertFiles); n++ {
		img, _, err := e.getConvertFileByArray(n)
		if err != nil {
			return err
		}
		page := strconv.Itoa(n + 1)
		name := filepath.Join(e.TargetFolder, e.PrefixFileName+page+e.ConvertFileFormat)
		err = e.createFile(img, name)
		if err != nil {
			return err
		}
	}
	return nil
}
