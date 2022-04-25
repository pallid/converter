package processor

import (
	"converter/pdf"
	"testing"
)

func TestEntity_Convert(t *testing.T) {

	conv := pdf.NewLocalConverter()

	// arrByte := make([]byte)
	// aArrByte := make([][]byte)

	first := &Entity{
		SourceFile:        "../source/test-folder1/pdf-test.pdf",
		TargetFolder:      "../target/test-folder1/pdf-test.pdf",
		ConvertFileFormat: DefaultFormat,
		PrefixFileName:    "",
		Converter:         conv,
		// bSourceFile:       arrByte,
		// bConvertFiles:     aArrByte,
	}

	tests := []struct {
		name    string
		e       *Entity
		wantErr bool
	}{
		{
			name:    "first",
			e:       first,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.Convert(); (err != nil) != tt.wantErr {
				t.Errorf("Entity.Convert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
