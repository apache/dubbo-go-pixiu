package urlPath

import (
	"strings"
)

// todo
func Split(path string) []string {
	return strings.Split(path, "/")
}

// {xxx} supported.
func IsWildcards(partOfPath string) bool {
	return strings.HasPrefix(partOfPath, "{") && strings.HasSuffix(partOfPath, "}")
}
