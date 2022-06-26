package path_parser

import (
	"regexp"
	"strings"
)

type Output struct {
	Protocols []string
	Protocol  string
	Port      int
	Resource  string
	User      string
	Pathname  string
	Hash      string
	Href      string
	Query     string
}

// IsSsh check if an input value is a ssh url or not
func IsSsh(input string) bool {
	protocols := Protocols(input)
	input = input[strings.Index(input, "://")+3:]

	for _, p := range protocols {
		if strings.HasPrefix(p, "ssh") || strings.HasPrefix(p, "rsync") {
			return true
		}
	}

	var re = regexp.MustCompile(`\.([a-zA-Z\\d]+):(\d+)`)
	return !re.Match([]byte(input)) && strings.Index(input, "@") < strings.Index(input, ":")
}

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
