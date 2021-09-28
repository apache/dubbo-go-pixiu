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

package http

import (
	"bytes"
	"io/ioutil"
	"mime"
	"net/http"
	"regexp"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

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

// ExtractRequestBody extract body of http request
func ExtractRequestBody(req *http.Request) []byte {
	isUpload := isUpload(req)
	if isUpload {
		return nil
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil
	}
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body
}

// isUpload determine whether it is a file upload request
func isUpload(req *http.Request) bool {
	var ret bool
	v := req.Header.Get("Content-Type")
	if v == "" {
		return ret
	}
	d, _, err := mime.ParseMediaType(v)
	if err == nil && d == "multipart/form-data" {
		ret = true
	}
	return ret
}
