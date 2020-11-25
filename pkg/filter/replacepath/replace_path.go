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

package replacepath

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	contexthttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter"
	"net/http"
	"net/url"
)


const (
	// ReplacedPathHeader is the default header to set the old path to.
	ReplacedPathHeader = "X-Replaced-Path"
	replacePathError = "replace path fail"
)

// replacePathFilter is a filter for host.
type replacePathFilter struct {
	path string
}

// New create replace path filter.
func New(path string) filter.Filter {
	return &replacePathFilter{path: path}
}

// // Do execute replacePathFilter filter logic.
func (f *replacePathFilter) Do() selfcontext.FilterFunc {
	return func(c selfcontext.Context) {
		f.doReplacePathFilter(c.(*contexthttp.HttpContext))
	}
}

func (f *replacePathFilter) doReplacePathFilter(ctx *contexthttp.HttpContext) {
	req := ctx.Request
	if req.URL.RawPath == "" {
		req.Header.Add(ReplacedPathHeader, req.URL.Path)
	} else {
		req.Header.Add(ReplacedPathHeader, req.URL.RawPath)
	}

	req.URL.RawPath = f.path
	var err error
	req.URL.Path, err = url.PathUnescape(req.URL.RawPath)
	if err != nil {
		ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		ctx.WriteWithStatus(http.StatusInternalServerError, []byte(replacePathError))
		ctx.Abort()
		return
	}

	req.RequestURI = req.URL.RequestURI()

	ctx.Next()
}
