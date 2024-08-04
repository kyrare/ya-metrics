package utils

import (
	"strings"
	"testing"
)

func TestHash(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Default",
			args: args{
				data: []byte("Default"),
			},
			want: "21b111cbfe6e8fca2d181c43f53ad548b22e38aca955b9824706a504b0a07a2d",
		},
		{
			name: "Empty String",
			args: args{
				data: []byte(""),
			},
			want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name: "Long String",
			args: args{
				data: []byte(strings.Repeat("a", 1000)),
			},
			want: "41edece42d63e8d9bf515a9ba6932e1c20cbc9f5a5d134645adb5db1b9737ea3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hash(tt.args.data); got != tt.want {
				t.Errorf("Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}
