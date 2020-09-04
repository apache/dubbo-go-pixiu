package http

import (
	"encoding/json"
	"net/http"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

// HttpContext
type HttpContext struct {
	*context.BaseContext
	HttpConnectionManager model.HttpConnectionManager
	FilterChains          []model.FilterChain
	Listener              *model.Listener
	api                   *client.Api

	Request   *http.Request
	writermem responseWriter
	Writer    ResponseWriter
}

func (hc *HttpContext) GetWritermem() {

}

func (hc *HttpContext) Next() {
	hc.Index++
	for hc.Index < int8(len(hc.Filters)) {
		hc.Filters[hc.Index](hc)
		hc.Index++
	}
}

func (hc *HttpContext) Reset() {
	hc.Writer = &hc.writermem
	hc.Filters = nil
	hc.Index = -1
}

func (hc *HttpContext) Status(code int) {
	hc.Writer.WriteHeader(code)
}

func (hc *HttpContext) StatusCode() int {
	return hc.Writer.Status()
}

func (hc *HttpContext) Write(b []byte) (int, error) {
	return hc.Writer.Write(b)
}

func (hc *HttpContext) WriteHeaderNow() {
	hc.writermem.WriteHeaderNow()
}

// WriteWithStatus status must set first
func (hc *HttpContext) WriteWithStatus(code int, b []byte) (int, error) {
	hc.Writer.WriteHeader(code)
	return hc.Writer.Write(b)
}

func (hc *HttpContext) AddHeader(k, v string) {
	hc.Writer.Header().Add(k, v)
}

func (hc *HttpContext) GetHeader(k string) string {
	return hc.Request.Header.Get(k)
}

func (hc *HttpContext) GetUrl() string {
	return hc.Request.URL.Path
}

func (hc *HttpContext) GetMethod() string {
	return hc.Request.Method
}

func (hc *HttpContext) Api(api *client.Api) {
	hc.api = api
}

func (hc *HttpContext) GetApi() *client.Api {
	return hc.api
}

func (hc *HttpContext) WriteFail() {
	hc.doWriteJSON(nil, http.StatusInternalServerError, nil)
}

func (hc *HttpContext) WriteErr(p interface{}) {
	hc.doWriteJSON(nil, http.StatusInternalServerError, p)
}

func (hc *HttpContext) WriteSuccess() {
	hc.doWriteJSON(nil, http.StatusOK, nil)
}

func (hc *HttpContext) WriteResponse(resp client.Response) {
	hc.doWriteJSON(nil, http.StatusOK, resp.Data)
}

func (hc *HttpContext) doWriteJSON(h map[string]string, code int, d interface{}) {
	if h == nil {
		h = make(map[string]string, 1)
	}
	h[constant.HeaderKeyContextType] = constant.HeaderValueJsonUtf8
	hc.doWrite(h, code, d)
}

func (hc *HttpContext) doWrite(h map[string]string, code int, d interface{}) {
	for k, v := range h {
		hc.Writer.Header().Set(k, v)
	}

	hc.Writer.WriteHeader(code)

	if d != nil {
		if b, err := json.Marshal(d); err != nil {
			hc.Writer.Write([]byte(err.Error()))
		} else {
			hc.Writer.Write(b)
		}
	}
}

// BuildFilters build filter, from config http_filters
func (hc *HttpContext) BuildFilters() {
	var ff []context.FilterFunc

	if hc.HttpConnectionManager.HttpFilters == nil {
		return
	}

	for _, v := range hc.HttpConnectionManager.HttpFilters {
		ff = append(ff, extension.GetMustFilterFunc(v.Name))
	}

	hc.AppendFilterFunc(ff...)
}

func (hc *HttpContext) ResetWritermen(w http.ResponseWriter) {
	hc.writermem.reset(w)
}
