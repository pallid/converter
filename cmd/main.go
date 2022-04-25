package main

import (
	"converter/pdf"
	"flag"

	"converter/services/catalog"
)

var sourceFolder string
var targerFolder string

const (
	defaultSource   = "./source"
	defaultTarget   = "./target"
	defaultFileName = "img%03d"
	defaultFormat   = ".jpg"
)

func main() {

	flag.StringVar(&sourceFolder, "src", defaultSource, "source folder")
	flag.StringVar(&targerFolder, "tar", defaultTarget, "target folder")

	flag.Parse()

	conv := pdf.NewLocalConverter()
	catalog.Start(sourceFolder, targerFolder, conv, 10)

}
