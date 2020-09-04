package context

import "math"

const abortIndex int8 = math.MaxInt8 / 2

// BaseContext
type BaseContext struct {
	Context
	Index   int8
	Filters FilterChain
}

func NewBaseContext() *BaseContext {
	return &BaseContext{Index: -1}
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *BaseContext) Next() {
	c.Index++
	for c.Index < int8(len(c.Filters)) {
		c.Filters[c.Index](c)
		c.Index++
	}
}

func (c *BaseContext) Abort() {
	c.Index = abortIndex
}

func (c *BaseContext) AppendFilterFunc(ff ...FilterFunc) {
	for _, v := range ff {
		c.Filters = append(c.Filters, v)
	}
}
