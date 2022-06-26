package path_parser

import (
	"strings"
)

// Protocols returns the protocols of an input url.
func Protocols(input string) []string {
	idx := strings.Index(input, "://")
	if idx == -1 {
		return []string{}
	}
	input = input[0:idx]
	splits := strings.Split(input, "+")

	var filtered []string

	for i := range splits {
		if splits[i] != "" {
			filtered = append(filtered, splits[i])
		}
	}

	return filtered
}
