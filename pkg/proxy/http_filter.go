package proxy

import (
	"regexp"
	"strings"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

func init() {
	AddFilterFunc(pkg.HttpDomainFilter, Domain())
}

// Domain
// https :authority
// http Host
func Domain() FilterFunc {
	return func(c Context) {
		if MatchDomainFilter(c.(*HttpContext)) {
			c.Next()
		}
	}
}

func MatchDomainFilter(c *HttpContext) bool {
	for _, v := range c.l.FilterChains {
		for _, d := range v.FilterChainMatch.Domains {
			if d == c.GetHeader("Host") {
				return true
			}
		}
	}

	return false
}

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

func HttpRouteActionMatch(c *HttpContext, ra model.RouteAction) bool {
	conf := GetBootstrap()

	if ra.Cluster == "" || !conf.ExistCluster(ra.Cluster) {
		return false
	}

	return true
}
