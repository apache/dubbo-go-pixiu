package proxy

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"math"
	"net/http"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

const abortIndex int8 = math.MaxInt8 / 2

type Context interface {
	Next()
	Abort()

	AppendFilterFunc(ff ...FilterFunc)

	Status(code int)
	StatusCode() int
	WriteWithStatus(int, []byte) (int, error)
	Write([]byte) (int, error)
	AddHeader(k, v string)
	GetHeader(k string) string
	GetUrl() string
	GetMethod() string

	BuildFilters()

	Api(api *Api)
	GetApi() *Api

	WriteErr(p interface{})
}

type BaseContext struct {
	Context
	index   int8
	filters _FilterChain
}

type FilterFunc func(Context)

type _FilterChain []FilterFunc

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *BaseContext) Next() {
	c.index++
	for c.index < int8(len(c.filters)) {
		c.filters[c.index](c)
		c.index++
	}
}

func (c *BaseContext) Abort() {
	c.index = abortIndex
}

func (c *BaseContext) AppendFilterFunc(ff ...FilterFunc) {
	for _, v := range ff {
		c.filters = append(c.filters, v)
	}
}

type HttpContext struct {
	*BaseContext
	httpConnectionManager model.HttpConnectionManager
	FilterChains          []model.FilterChain
	l                     *model.Listener
	api                   *Api

	r         *http.Request
	writermem responseWriter
	Writer    ResponseWriter
}

func (hc *HttpContext) Next() {
	hc.index++
	for hc.index < int8(len(hc.filters)) {
		hc.filters[hc.index](hc)
		hc.index++
	}
}

func (hc *HttpContext) reset() {
	hc.Writer = &hc.writermem
	hc.filters = nil
	hc.index = -1
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

func (hc *HttpContext) WriteWithStatus(code int, b []byte) (int, error) {
	hc.Writer.WriteHeader(code)
	return hc.Writer.Write(b)
}

func (hc *HttpContext) AddHeader(k, v string) {
	hc.Writer.Header().Add(k, v)
}

func (hc *HttpContext) GetHeader(k string) string {
	return hc.r.Header.Get(k)
}

func (hc *HttpContext) GetUrl() string {
	return hc.r.URL.Path
}

func (hc *HttpContext) GetMethod() string {
	return hc.r.Method
}

func (hc *HttpContext) Api(api *Api) {
	hc.api = api
}

func (hc *HttpContext) GetApi() *Api {
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

func (hc *HttpContext) WriteResponse(resp Response) {
	hc.doWriteJSON(nil, http.StatusOK, resp.data)
}

func (hc *HttpContext) doWriteJSON(h map[string]string, code int, d interface{}) {
	if h == nil {
		h = make(map[string]string, 1)
	}
	h[pkg.HeaderKeyContextType] = pkg.HeaderValueJsonUtf8
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

func (hc *HttpContext) BuildFilters() {
	var ff []FilterFunc

	for _, v := range hc.httpConnectionManager.HttpFilters {
		ff = append(ff, GetMustFilterFunc(v.Name))
	}

	hc.AppendFilterFunc(ff...)
}

func GetBytes(k interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(k)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
