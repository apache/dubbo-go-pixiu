package triple

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	proxymeta "github.com/mercari/grpc-http-proxy/metadata"
	"github.com/mercari/grpc-http-proxy/proxy"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
)

// InitDefaultTripleClient init default dubbo client
func InitDefaultTripleClient() {
	tripleClient = NewTripleClient()
}

var tripleClient *Client
var clientOnc sync.Once

// NewTripleClient create dubbo client
func NewTripleClient() *Client {
	clientOnc.Do(func() {
		tripleClient = &Client{}
	})
	return tripleClient
}

// Client client to generic invoke dubbo
type Client struct {
}

func (tc *Client) Apply() error {
	panic("implement me")
}

func (tc *Client) MapParams(req *client.Request) (reqData interface{}, err error) {
	panic("implement me")
}

// Close clear GenericServicePool.
func (dc *Client) Close() error {
	return nil
}

// SingletonTripleClient singleton dubbo clent
func SingletonTripleClient() *Client {
	return NewTripleClient()
}

// Call invoke service
func (dc *Client) Call(req *client.Request) (res interface{}, err error) {
	address := strings.Split(req.API.IntegrationRequest.HTTPBackendConfig.URL, ":")
	p := proxy.NewProxy()
	targetURL := &url.URL{
		Scheme: address[0],
		Opaque: address[1],
	}
	if err := p.Connect(context.Background(), targetURL); err != nil {
		panic(err)
	}
	meta := make(map[string][]string)
	reqData, _ := ioutil.ReadAll(req.IngressRequest.Body)
	call, err := p.Call(context.Background(), req.API.Method.IntegrationRequest.Interface, req.API.Method.IntegrationRequest.Method, reqData, (*proxymeta.Metadata)(&meta))
	if err != nil {
		panic(err)
	}
	return string(call), nil
}
