package filter

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

//FilterChain
type FilterChain interface {
	AppendDecodeFilters(f ...HttpDecodeFilter)
	AppendEncodeFilters(f ...HttpEncodeFilter)
	OnDecode(ctx *http.HttpContext)
	OnEncode(ctx *http.HttpContext)
}

type defaultFilterChain struct {
	decodeFilters      []HttpDecodeFilter
	decodeFiltersIndex int

	encodeFilters      []HttpEncodeFilter
	encodeFiltersIndex int
}

func newDefaultFilterChain() FilterChain {
	return &defaultFilterChain{
		decodeFilters:      []HttpDecodeFilter{},
		decodeFiltersIndex: 0,
		encodeFilters:      []HttpEncodeFilter{},
		encodeFiltersIndex: 0,
	}
}

func (c *defaultFilterChain) AppendDecodeFilters(f ...HttpDecodeFilter) {
	c.decodeFilters = append(c.decodeFilters, f...)
}

// AppendEncodeFilters append encode filters in reverse order
func (c *defaultFilterChain) AppendEncodeFilters(f ...HttpEncodeFilter) {
	for i := len(f) - 1; i >= 0; i-- {
		c.encodeFilters = append([]HttpEncodeFilter{f[i]}, c.encodeFilters...)
	}
}

func (c *defaultFilterChain) OnDecode(ctx *http.HttpContext) {
	for ; c.decodeFiltersIndex < len(c.decodeFilters); c.decodeFiltersIndex++ {
		filterStatus := c.decodeFilters[c.decodeFiltersIndex].Decode(ctx)

		switch filterStatus {
		case Continue:
			continue
		case Stop:
			return
		}
	}
}

func (c *defaultFilterChain) OnEncode(ctx *http.HttpContext) {
	for ; c.encodeFiltersIndex < len(c.encodeFilters); c.encodeFiltersIndex++ {
		filterStatus := c.encodeFilters[c.encodeFiltersIndex].Encode(ctx)

		switch filterStatus {
		case Continue:
			continue
		case Stop:
			return
		}
	}
}
