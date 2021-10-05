package model

import (
	"net/http"
	"net/url"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestRouterPathMatch(t *testing.T) {
	rm := &RouterMatch{
		Prefix:  "/user/",
		Path:    "/user2/test",
		Regex:   "/user3/test",
		Methods: []string{"POST"},
		Headers: []HeaderMatcher{
			{
				Name: "content-type",
				Values: []string{
					"grpc",
					"json",
				},
			},
		},
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user/add",
			},
		}
		assert.True(t, rm.matchPath(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user2/add",
			},
		}
		assert.False(t, rm.matchPath(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user2/test",
			},
		}
		assert.True(t, rm.matchPath(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
		}
		assert.True(t, rm.matchPath(req))
	}
}

func TestRouterMethodMatch(t *testing.T) {
	rm := &RouterMatch{
		Prefix:  "/user",
		Methods: []string{"POST", "PUT"},
		Headers: []HeaderMatcher{
			{
				Name: "content-type",
				Values: []string{
					"grpc",
					"json",
				},
			},
		},
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "POST",
		}
		assert.True(t, rm.matchMethod(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "PUT",
		}
		assert.True(t, rm.matchMethod(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "GET",
		}
		assert.False(t, rm.matchMethod(req))
	}
}

func TestRouterHeaderMatch(t *testing.T) {
	rm := &RouterMatch{
		Prefix:  "/user",
		Path:    "/user2/test",
		Regex:   "/user3/test",
		Methods: []string{"POST", "PUT"},
		Headers: []HeaderMatcher{
			{
				Name: "content-type",
				Values: []string{
					"grpc",
				},
			},
		},
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "POST",
			Header: map[string][]string{
				"Content-Type": []string{
					"grpc",
				},
			},
		}
		assert.True(t, rm.matchHeader(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "POST",
		}
		assert.False(t, rm.matchHeader(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "POST",
			Header: map[string][]string{
				"Content-Type": []string{
					"grpc2",
				},
			},
		}
		assert.False(t, rm.matchHeader(req))
	}
}

func TestRouterMatch(t *testing.T) {
	rm := RouterMatch{
		Prefix:  "/user/",
		Path:    "/user2/test",
		Regex:   "/user3/test",
		Methods: []string{"POST"},
		Headers: []HeaderMatcher{
			{
				Name: "content-type",
				Values: []string{
					"grpc",
					"json",
				},
			},
		},
	}

	router := &Router{
		Match: rm,
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "POST",
			Header: map[string][]string{
				"Content-Type": []string{
					"grpc",
				},
			},
		}
		assert.True(t, router.MatchRouter(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/us3/test",
			},
			Method: "POST",
			Header: map[string][]string{
				"Content-Type": []string{
					"grpc",
				},
			},
		}
		assert.False(t, router.MatchRouter(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "POST11",
			Header: map[string][]string{
				"Content-Type": []string{
					"grpc",
				},
			},
		}
		assert.False(t, router.MatchRouter(req))
	}

	{
		req := &http.Request{
			URL: &url.URL{
				Path: "/user3/test",
			},
			Method: "POST",
			Header: map[string][]string{
				"Content-Type": []string{
					"grpc1111",
				},
			},
		}
		assert.False(t, router.MatchRouter(req))
	}
}
