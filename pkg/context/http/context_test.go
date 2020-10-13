package http

import (
	"testing"
)

func TestBuildContext(t *testing.T) {
	ctx := HttpContext{}
	ctx.BuildFilters()
}
