package path_parser

import (
	"github.com/stretchr/testify/require"
	"strconv"
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
			require.Equal(t, Protocols(test.input), test.expected)
		})
	}
}

func TestIsSsh(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		// Secure Shell Transport Protocol (SSH)
		{"ssh://user@host.xz:port/path/to/repo.git/", true},
		{"ssh://user@host.xz/path/to/repo.git/", true},
		{"ssh://host.xz:port/path/to/repo.git/", true},
		{"ssh://host.xz/path/to/repo.git/", true},
		{"ssh://user@host.xz/path/to/repo.git/", true},
		{"ssh://host.xz/path/to/repo.git/", true},
		{"ssh://user@host.xz/~user/path/to/repo.git/", true},
		{"ssh://host.xz/~user/path/to/repo.git/", true},
		{"ssh://user@host.xz/~/path/to/repo.git", true},
		{"ssh://host.xz/~/path/to/repo.git", true},
		{"user@host.xz:/path/to/repo.git/", true},
		{"user@host.xz:~user/path/to/repo.git/", true},
		{"user@host.xz:path/to/repo.git", true},
		{"host.xz:/path/to/repo.git/", true},
		{"host.xz:path/to/repo.git", true},
		{"host.xz:~user/path/to/repo.git/", true},
		{"rsync://host.xz/path/to/repo.git/", true},

		// Git Transport Protocol
		{"git://host.xz/path/to/repo.git/", false},
		{"git://host.xz/~user/path/to/repo.git/", false},

		// HTTP/S Transport Protocol
		{"http://host.xz/path/to/repo.git/", false},
		{"https://host.xz/path/to/repo.git/", false},
		{"http://host.xz:8000/path/to/repo.git/", false},
		{"https://host.xz:8000/path/to/repo.git/", false},
		// Local (Filesystem) Transport Protocol
		{"/path/to/repo.git/", false},
		{"path/to/repo.git/", false},
		{"~/path/to/repo.git", false},
		{"file:///path/to/repo.git/", false},
		{"file://~/path/to/repo.git/", false},
	}

	for n, test := range tests {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			require.Equal(t, test.expected, IsSsh(test.input))
		})
	}
}
