package triple

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/registry"
	_ "dubbo.apache.org/dubbo-go/v3/registry/nacos"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	proxymeta "github.com/mercari/grpc-http-proxy/metadata"
	"github.com/mercari/grpc-http-proxy/proxy"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"
	"time"
)
// InitDefaultTripleClient init default dubbo client
func InitDefaultTripleClient(dpc *TripleProxyConfig) {
	tripleClient = NewTripleClient()
	tripleClient.setConfig(dpc)
	if err := tripleClient.Apply(); err != nil {
		logger.Warnf("dubbo client apply error %s", err)
	}
}

var tripleClient *Client
var clientOnc sync.Once
// NewTripleClient create dubbo client
func NewTripleClient() *Client {
	clientOnc.Do(func() {
		tripleClient = &Client{
				lock:               sync.RWMutex{},
				zkNotifyListenerInstance: newZKNotifyListener(),
			}
		})
	return tripleClient
}

// SetConfig set config
func (tc *Client) setConfig(dpc *TripleProxyConfig) {
	tc.tripleProxyConfig = dpc
}

var once = sync.Once{}
var zkNotifyListenerInstance *zkNotifyListener

func newZKNotifyListener() *zkNotifyListener{
	once.Do(func() {
		zkNotifyListenerInstance = &zkNotifyListener{
			serviceDiscMap: make(map[string]string),
		}
	})
	return zkNotifyListenerInstance
}

type zkNotifyListener struct {
	lock sync.Mutex
	serviceDiscMap map[string]string
}

func (l *Client) getAddr(interfaceName string) (string,bool) {
	if addr, ok := l.zkNotifyListenerInstance.serviceDiscMap[interfaceName]; ok{
		return addr, true
	}
	return "", false
}

func (l *zkNotifyListener) Notify(e* registry.ServiceEvent){
	l.lock.Lock()
	defer l.lock.Unlock()
	l.serviceDiscMap[ e.Service.GetParam(constant.InterfaceKey, "")] = e.Service.Location
	fmt.Println("get event = ",e)
}
func (l *zkNotifyListener) NotifyAll(e[]*registry.ServiceEvent,f func()){
	fmt.Println("get event = ",e)
	f()
}

func (tc*Client)subscribe(interfaceName string)error{
	if _, ok := tc.subscribedKey.Load(interfaceName); ok{
		return nil
	}

	url, err := common.NewURL("tri://localhost:20000")
	if err != nil{
		return err
	}

	url.SetParam(constant.InterfaceKey, interfaceName)
	url.SetParam(constant.GroupKey, "")
	url.SetParam(constant.VersionKey, "")
	go func() {
		if err := tc.regInstance.Subscribe(url, tc.zkNotifyListenerInstance); err != nil{
			return
		}
	}()
	tc.subscribedKey.Store(interfaceName, true)
	return nil
}

// Apply init dubbo, config mapping can do here
func (tc *Client) Apply() error {
	pixiuRegConfig := tc.tripleProxyConfig.Registries
	regConfig := config.NewRegistryConfigBuilder().
		SetProtocol(pixiuRegConfig.Protocol).
		SetAddress(pixiuRegConfig.Address).
		SetTimeout(pixiuRegConfig.Timeout).
		SetUsername(pixiuRegConfig.Username).
		SetPassword(pixiuRegConfig.Password).
		Build()
	regInstance, err := regConfig.GetInstance(common.CONSUMER)
	if err != nil{
		return err
	}
	tc.regInstance = regInstance
	return nil
}


// Client client to generic invoke dubbo
type Client struct {
	lock               sync.RWMutex
	tripleProxyConfig   *TripleProxyConfig
	zkNotifyListenerInstance *zkNotifyListener
	regInstance registry.Registry
	subscribedKey sync.Map
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
	after := time.After(time.Second*3)
	var addr string
LOOP:
	for {
		select {
		case <- after:
			return nil, errors.Errorf("no provider named %s", req.API.Interface)
		default:
			var ok bool
			if addr, ok = dc.getAddr(req.API.Method.IntegrationRequest.Interface); ok{
				break LOOP
			}else{
				if err := dc.subscribe(req.API.Method.IntegrationRequest.Interface); err != nil{
					return nil, err
				}
			}
		}
		time.Sleep(time.Millisecond*100)
	}
	addrs := strings.Split(addr, ":")
	p := proxy.NewProxy()
	targetURL := &url.URL{
		Scheme: addrs[0],
		Opaque: addrs[1],
	}
	if err := p.Connect(context.Background(), targetURL); err != nil{
		panic(err)
	}
	meta := make(map[string][]string)
	reqData, _ := ioutil.ReadAll( req.IngressRequest.Body)
	call, err := p.Call(context.Background(), req.API.Method.IntegrationRequest.Interface, req.API.Method.IntegrationRequest.Method, reqData, (*proxymeta.Metadata)(&meta))
	if err != nil{
		panic(err)
	}
	return string(call), nil
}
