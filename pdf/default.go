package pdf

type defaultConverter struct{}

func NewLocalConverter() Converter {
	return defaultConverter{}
}

func (defaultConverter) Convert(file []byte, format imageFormat) ([][]byte, error) {
	converter := newMuPdfConverter()
	images, err := converter.Convert(file, ConversionOptions{OutputFormat: format})
	if err != nil {
		return nil, err
	}
	return images, nil
}
