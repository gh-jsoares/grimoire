package clipboard

import (
	"encoding/base64"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"a", "YQ=="},
		{"ab", "YWI="},
		{"abc", "YWJj"},
		{"hello", "aGVsbG8="},
		{"Hello, World!", "SGVsbG8sIFdvcmxkIQ=="},
		{"x", "eA=="},
		{"xy", "eHk="},
		{"xyz", "eHl6"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := encode(tt.input)
			if got != tt.want {
				t.Errorf("encode(%q) = %q, want %q", tt.input, got, tt.want)
			}

			// Cross-check against stdlib
			stdlib := base64.StdEncoding.EncodeToString([]byte(tt.input))
			if got != stdlib {
				t.Errorf("encode(%q) = %q, stdlib = %q", tt.input, got, stdlib)
			}
		})
	}
}

func TestEncode_BinaryData(t *testing.T) {
	input := "\x00\x01\x02\xff"
	got := encode(input)
	want := base64.StdEncoding.EncodeToString([]byte(input))
	if got != want {
		t.Errorf("encode(binary) = %q, want %q", got, want)
	}
}
