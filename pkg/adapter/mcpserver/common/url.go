package common

import (
	"net"
	"net/url"
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

	// If it looks like a URL with scheme, try url.Parse first
	if i := strings.Index(addr, "://"); i >= 0 {
		u, err := url.Parse(addr)
		if err == nil && u.Host != "" {
			return splitHostPortCompat(u.Host)
		}
		// fallback by stripping scheme and path
		addr = addr[i+3:]
	}

	// Strip path/query if present
	if j := strings.IndexAny(addr, "/?\n\r\t "); j >= 0 {
		addr = addr[:j]
	}

	return splitHostPortCompat(addr)
}

// splitHostPortCompat splits host:port using net.SplitHostPort if possible, with a fallback.
func splitHostPortCompat(hostport string) (string, int) {
	host, portStr, err := net.SplitHostPort(hostport)
	if err != nil {
		// try naive split on last ':' for simple IPv4/hostname cases
		idx := strings.LastIndex(hostport, ":")
		if idx <= 0 || idx == len(hostport)-1 {
			return "", 0
		}
		host = hostport[:idx]
		portStr = hostport[idx+1:]
	}
	// parse int port
	port := 0
	for i := 0; i < len(portStr); i++ {
		c := portStr[i]
		if c < '0' || c > '9' {
			return "", 0
		}
		port = port*10 + int(c-'0')
	}
	if port <= 0 {
		return "", 0
	}
	return strings.TrimSpace(host), port
}
