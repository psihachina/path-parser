package path_parser

import (
	"reflect"
	"testing"
)

func TestProtocols(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{name: "multiple protocols", input: "git+ssh://git@some-host.com/and-the-path/name", expected: []string{"git", "ssh"}},
		{name: "no protocols", input: "//foo.com/bar.js", expected: []string{}},
		{name: "one protocol", input: "ssh://git@some-host.com/and-the-path/name", expected: []string{"ssh"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !reflect.DeepEqual(Protocols(test.input), test.expected) {
				t.Fail()
			}
		})
	}
}
