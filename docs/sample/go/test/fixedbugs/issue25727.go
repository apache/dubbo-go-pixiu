// errorcheck

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "net/http"

var s = http.Server{}
var _ = s.doneChan // ERROR "s.doneChan undefined .cannot refer to unexported field or method doneChan.$"
var _ = s.DoneChan // ERROR "s.DoneChan undefined .type http.Server has no field or method DoneChan.$"
var _ = http.Server{tlsConfig: nil} // ERROR "unknown field 'tlsConfig' in struct literal.+ .but does have TLSConfig.$"
var _ = http.Server{DoneChan: nil} // ERROR "unknown field 'DoneChan' in struct literal of type http.Server$"

type foo struct {
    bar int
}

var _ = &foo{bAr: 10} // ERROR "unknown field 'bAr' in struct literal.+ .but does have bar.$"
