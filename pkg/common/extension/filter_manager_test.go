package extension

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	pxCtx "github.com/apache/dubbo-go-pixiu/pkg/context"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func init() {
	RegisterFilterFactory(DEMO, newDemoFilter)
}

const DEMO = "dgp.filters.demo"

type DemoFilter struct {
	conf *Config
}

type Config struct {
	Foo string `json:"foo,omitempty" yaml:"foo,omitempty"`
	Bar string `json:"bar,omitempty" yaml:"bar,omitempty"`
}

func newDemoFilter() filter.Factory {
	return &DemoFilter{conf: &Config{Foo: "default foo", Bar: "default bar"}}
}

func (d *DemoFilter) Config() interface{} {
	return d.conf
}

func (d *DemoFilter) Apply() (filter.Filter, error) {
	//init every sth before return the func
	c := d.conf
	str := fmt.Sprintf("%s is drinking in the %s", c.Foo, c.Bar)
	//return the filter func
	return func(ctx context.Context) {
		logger.Info(str)
	}, nil
}

func TestApply(t *testing.T) {
	fm := NewFilterManager()

	conf := map[string]interface{}{}
	conf["foo"] = "Cat"
	conf["bar"] = "The Walnut"
	apply, err := fm.Apply(DEMO, conf)
	assert.Nil(t, err)

	baseContext := pxCtx.NewBaseContext()
	apply(baseContext)
}
