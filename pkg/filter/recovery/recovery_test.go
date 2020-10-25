package recovery

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

import (
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	selfhttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

func TestRecovery(t *testing.T) {
	c := MockHttpContext(func(c selfcontext.Context) {
		time.Sleep(time.Millisecond * 100)
		// panic
		var m map[string]string
		m["name"] = "1"
	})
	c.Next()
	// print
	// 500
	// "assignment to entry in nil map"
}

func MockHttpContext(fc ...selfcontext.FilterFunc) *selfhttp.HttpContext {
	result := &selfhttp.HttpContext{
		BaseContext: &selfcontext.BaseContext{
			Index: -1,
			Ctx:   context.Background(),
		},
	}

	w := MockWriter{header: map[string][]string{}}
	result.ResetWritermen(&w)
	result.Reset()

	result.Filters = append(result.Filters, NewRecoveryFilter().Do())
	for i := range fc {
		result.Filters = append(result.Filters, fc[i])
	}

	return result
}

type MockWriter struct {
	header http.Header
}

func (w *MockWriter) Header() http.Header {
	return w.header
}

func (w *MockWriter) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return -1, nil
}

func (w *MockWriter) WriteHeader(statusCode int) {
	fmt.Println(statusCode)
}
