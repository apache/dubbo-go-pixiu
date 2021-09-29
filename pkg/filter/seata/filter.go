package seata

import (
	http2 "net/http"
	"strings"

	"github.com/opentrx/seata-golang/v2/pkg/apis"
	"github.com/opentrx/seata-golang/v2/pkg/util/runtime"
	"google.golang.org/grpc"

	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	Kind = "dgp.filter.http.seata"

	SEATA    = "seata"
	XID      = "x_seata_xid"
	BranchID = "x_seata_branch_id"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// MetricFilter is http Filter plugin.
	Plugin struct {
	}
	// Filter is http Filter instance
	Filter struct {
		conf              *Seata
		transactionInfos  map[string]*TransactionInfo
		tccResources      map[string]*TCCResource
		transactionClient apis.TransactionManagerServiceClient
		resourceClient    apis.ResourceManagerServiceClient
		branchMessages    chan *apis.BranchMessage
	}
)

func (ap *Plugin) Kind() string {
	return Kind
}

func (ap *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{
		conf:             &Seata{},
		transactionInfos: make(map[string]*TransactionInfo),
		tccResources:     make(map[string]*TCCResource),
		branchMessages:   make(chan *apis.BranchMessage),
	}, nil
}

func (f *Filter) Config() interface{} {
	return f.conf
}

func (f *Filter) Apply() error {
	conn, err := grpc.Dial(f.conf.ServerAddressing,
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(f.conf.GetClientParameters()))
	if err != nil {
		return err
	}

	f.transactionClient = apis.NewTransactionManagerServiceClient(conn)
	f.resourceClient = apis.NewResourceManagerServiceClient(conn)

	runtime.GoWithRecover(func() {
		f.branchCommunicate()
	}, nil)

	for _, ti := range f.conf.TransactionInfos {
		f.transactionInfos[ti.RequestPath] = ti
	}

	for _, r := range f.conf.TCCResources {
		f.tccResources[r.PrepareRequestPath] = r
	}
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *http.HttpContext) {
	path := ctx.Request.URL.Path
	method := ctx.Request.Method

	if method != http2.MethodPost {
		ctx.Next()
		return
	}

	transactionInfo, found := f.transactionInfos[strings.ToLower(path)]
	if found {
		result := f.handleHttp1GlobalBegin(ctx, transactionInfo)
		ctx.Next()
		if result {
			f.handleHttp1GlobalEnd(ctx)
		}
	}

	tccResource, exists := f.tccResources[strings.ToLower(path)]
	if exists {
		result := f.handleHttp1BranchRegister(ctx, tccResource)
		ctx.Next()
		if result {
			f.handleHttp1BranchEnd(ctx)
		}
	}
}
