package catalog

import "testing"

func Test_isSupportDocument(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "pdf ext",
			args: args{
				fileName: "test.pdf",
			},
			want: true,
		},
		{
			name: "doc ext",
			args: args{
				fileName: "test.doc",
			},
			want: false,
		},
		{
			name: "fake ext",
			args: args{
				fileName: "test.fake",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSupportDocument(tt.args.fileName); got != tt.want {
				t.Errorf("isSupportDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}
