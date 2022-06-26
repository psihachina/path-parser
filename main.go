package path_parser

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type Output struct {
	Protocols []string
	Protocol  string
	Port      *int
	Resource  string
	User      string
	Pathname  string
	Hash      string
	Search    string
	Href      string
	Query     url.Values
}

// ParsePath parse paths (local paths, urls: ssh/git/etc)
func ParsePath(input string) Output {
	input = strings.TrimSpace(input)
	var re = regexp.MustCompile(`\r\n|\r`)
	input = string(re.ReplaceAll([]byte(input), []byte("")))
	var output = Output{
		Protocols: Protocols(input),
		Protocol:  "",
		Port:      nil,
		Resource:  "",
		User:      "",
		Pathname:  "",
		Hash:      "",
		Search:    "",
		Href:      input,
		Query:     nil,
	}

	if strings.HasPrefix(input, ".") {
		if strings.HasPrefix(input, "./") {
			input = input[2:]
		}
		output.Pathname = input
		output.Protocol = "file"
	}

	protocolIdx := strings.Index(input, "://")

	fc := input[1]
	if len(output.Protocol) == 0 {
		if len(output.Protocols) > 0 {
			output.Protocol = output.Protocols[0]
		}

		if len(output.Protocol) == 0 {
			if IsSsh(input) {
				output.Protocol = "ssh"
			} else if fc == '/' || fc == '~' {
				input = input[2:]
				output.Protocol = "file"
			} else {
				output.Protocol = "file"
			}
		}
	}

	if protocolIdx != -1 {
		input = input[protocolIdx+3:]
	}

	re = regexp.MustCompile(`\/|\\`)
	parts := re.Split(input, -1)
	if output.Protocol != "file" {
		output.Resource = parts[0]
		parts = parts[1:]
	} else {
		output.Resource = ""
	}

	splits := strings.Split(output.Resource, "@")
	if len(splits) == 2 {
		output.User = splits[0]
		output.Resource = splits[1]
	}

	splits = strings.Split(output.Resource, ":")
	if len(splits) == 2 {
		output.Resource = splits[0]
		port := splits[1]
		if len(port) != 0 {
			p, err := strconv.Atoi(port)
			if err != nil || p == 0 {
				parts = append(parts, parts[len(parts)-1])
				parts[0] = port
			} else {
				output.Port = &p
			}
		}
	}

	var filtered []string
	for i := range parts {
		if parts[i] != "" {
			filtered = append(filtered, parts[i])
		}
	}

	parts = filtered

	if output.Protocol == "file" {
		output.Pathname = output.Href
	} else if len(output.Pathname) == 0 {
		pathname := ""
		if output.Protocol != "file" || output.Href[0] == '/' {
			pathname = "/"
		}

		output.Pathname = pathname + strings.Join(parts, "/")
	}

	splits = strings.Split(output.Pathname, "#")
	if len(splits) == 2 {
		output.Pathname = splits[0]
		output.Hash = splits[1]
	}

	splits = strings.Split(output.Pathname, "?")
	if len(splits) == 2 {
		output.Pathname = splits[0]
		output.Search = splits[1]
	}

	re = regexp.MustCompile(`\/$`)

	q, err := url.ParseQuery(output.Search)
	if err != nil {
		panic(err)
	}

	output.Query = q
	output.Href = re.ReplaceAllString(output.Href, "")
	output.Pathname = re.ReplaceAllString(output.Pathname, "")

	return output
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

	var re = regexp.MustCompile(`\.([a-zA-Z\\d]+):(\d+)\/`)
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
