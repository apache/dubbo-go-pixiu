package grpc

import (
	"net/http"
)

import (
	"golang.org/x/net/http2"
)

type HttpForwarder struct {
	transport *http2.Transport
}

func (hf *HttpForwarder) Forward(r *http.Request) (*http.Response, error) {
	return hf.transport.RoundTrip(r)
}
