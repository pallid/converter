package pdf

type imageFormat int32

const (
	ImageFormatUNKNOWN = imageFormat(0)
	ImageFormatJPEG    = imageFormat(1)
	ImageFormatPNG     = imageFormat(2)
	ImageFormatSVG     = imageFormat(3)
)

type Converter interface {
	Convert(file []byte, format imageFormat) ([][]byte, error)
}
