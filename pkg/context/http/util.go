package http

import (
	"regexp"
	"strings"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

// HttpHeaderMatch
func HttpHeaderMatch(c *HttpContext, hm model.HeaderMatcher) bool {
	if hm.Name == "" {
		return true
	}

	if hm.Value == "" {
		if c.GetHeader(hm.Name) == "" {
			return true
		}
	} else {
		if hm.Regex {
			// TODO
			return true
		} else {
			if c.GetHeader(hm.Name) == hm.Value {
				return true
			}
		}
	}

	return false
}

// HttpRouteMatch
func HttpRouteMatch(c *HttpContext, rm model.RouterMatch) bool {
	if rm.Prefix != "" {
		if !strings.HasPrefix(c.GetUrl(), rm.Path) {
			return false
		}
	}

	if rm.Path != "" {
		if c.GetUrl() != rm.Path {
			return false
		}
	}

	if rm.Regex != "" {
		if !regexp.MustCompile(rm.Regex).MatchString(c.GetUrl()) {
			return false
		}
	}

	return true
}

// HttpRouteActionMatch
func HttpRouteActionMatch(c *HttpContext, ra model.RouteAction) bool {
	conf := config.GetBootstrap()

	if ra.Cluster == "" || !conf.ExistCluster(ra.Cluster) {
		return false
	}

	return true
}
