package grpc

import (
	"golang.org/x/net/http2"
	"net/http"
)

type HttpForwarder struct {
	transport *http2.Transport
}

func (hf *HttpForwarder) Forward(r *http.Request) (*http.Response, error) {
	return hf.transport.RoundTrip(r)
}