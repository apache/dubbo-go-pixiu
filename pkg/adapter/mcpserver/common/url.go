package common

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

// ParseHostPortFromURL extracts host and port from a raw URL or host:port string.
// It supports forms like:
// - http://host:port/path?query
// - https://host:port
// - host:port
// Returns ("", 0) if parsing fails or no valid port is present.
func ParseHostPortFromURL(raw string) (string, int) {
	if raw == "" {
		return "", 0
	}

	addr := strings.TrimSpace(raw)

	if strings.Contains(addr, "://") {
		u, err := url.Parse(addr)
		if err == nil && u.Host != "" {
			return splitHostPort(u.Host)
		}
	}

	return splitHostPort(addr)
}

func splitHostPort(hostport string) (string, int) {
	host, portStr, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", 0
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		return "", 0
	}

	return strings.TrimSpace(host), port
}
