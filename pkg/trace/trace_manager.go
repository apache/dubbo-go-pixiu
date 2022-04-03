package trace

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type traceGenerator func(Trace) Trace

var TraceFactory map[ProtocolName]traceGenerator

func init() {
	TraceFactory = make(map[ProtocolName]traceGenerator)
	TraceFactory[HTTP] = NewHTTPTracer
}

type TraceDriverManager struct {
	driver    *TraceDriver
	bootstrap *model.Bootstrap
}

func CreateDefaultTraceDriverManager(bs *model.Bootstrap) *TraceDriverManager {
	manager := &TraceDriverManager{
		bootstrap: bs,
	}
	manager.driver = NewTraceDriver()
	manager.driver.Init(bs)
	return manager
}

func (manager *TraceDriverManager) GetDriver() *TraceDriver {
	return manager.driver
}

func (manager *TraceDriverManager) Close() {
	driver := manager.driver
	_ = driver.Tp.Shutdown(driver.context)
}
