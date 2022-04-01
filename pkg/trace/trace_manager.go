package trace

import (
	"errors"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/apache/dubbo-go-pixiu/pkg/trace/protocol"
	"strings"
	"sync/atomic"
)

type traceGenerator func(protocol.Trace) protocol.Trace

var traceFactory map[protocolName]traceGenerator

func init() {
	traceFactory = make(map[protocolName]traceGenerator)
	traceFactory[HTTP] = protocol.NewHTTPTracer
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

func NewTracer(name protocolName) protocol.Trace {
	driver := server.GetTraceDriverManager().driver
	holder, ok := driver.holders[name]
	if !ok {
		holder = &Holder{
			Tracers: make(map[string]protocol.Trace),
		}
		holder.Id = 0
		driver.holders[name] = holder
	}
	// tarceId的生成，并且在协议接口唯一
	builder := strings.Builder{}
	builder.WriteString(string(name))
	builder.WriteString("-" + string(holder.Id))

	traceId := builder.String()
	tmp := driver.tp.Tracer(traceId)
	tracer := &protocol.Tracer{
		Id: traceId,
		T:  tmp,
		H:  holder,
	}

	holder.Tracers[traceId] = traceFactory[HTTP](tracer)

	atomic.AddUint64(&holder.Id, 1)
	return holder.Tracers[traceId]
}

func GetTracer(name protocolName, tracerId string) (protocol.Trace, error) {
	driver := server.GetTraceDriverManager().driver
	holder, ok := driver.holders[name]
	if !ok {
		return nil, errors.New("can not find any tracer, please call NewTracer first")
	} else if _, ok = holder.Tracers[tracerId]; !ok {
		return nil, errors.New(fmt.Sprintf("can not find tracer %s with protocol %s", tracerId, name))
	} else {
		return holder.Tracers[tracerId], nil
	}
}

func (manager *TraceDriverManager) Close() {
	driver := manager.driver
	_ = driver.tp.Shutdown(driver.context)
}
