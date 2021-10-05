/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
