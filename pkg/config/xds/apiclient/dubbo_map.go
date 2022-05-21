package apiclient

import (
	"encoding/json"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	authorizationHeader    = "Authorization"
	istiodTokenPrefix      = "Bearer "
	DefaultIstiodTokenPath = "/var/run/secrets/token/istio-token"
)

// ADSZResponse is body got from istiod:8080/debug/adsz
type (
	ADSZResponse struct {
		Clients []ADSZClient `json:"clients"`
	}

	ADSZClient struct {
		Metadata map[string]interface{} `json:"metadata"`
		Address  string                 `json:"address"`
	}

	InterfaceMapHandlerImpl struct {
		istioTokenPath string
		istioDebugAddr string
	}

	DubboDetail struct {
		Service   string
		DSN       string
		Endpoints []string
	}
)

func (a *ADSZResponse) GetMap() map[string]*DubboDetail {
	result := make(map[string]*DubboDetail)
	for _, c := range a.Clients {
		// todo assert failed panic
		if dubboGoStr, ok := c.Metadata["LABELS"].(map[string]interface{})["DUBBO_GO"].(string); ok {
			resultMap := make(map[string]*string)
			_ = json.Unmarshal([]byte(dubboGoStr), &resultMap)
			for k, v := range resultMap {
				detail, exist := result[k]
				if !exist {
					detail = &DubboDetail{
						Endpoints: make([]string, 0),
					}
					detail.DSN = *v
					detail.Service = k
					result[k] = detail
				}
				address := strings.Split(c.Address, ":")
				detail.Endpoints = append(detail.Endpoints, strings.Join(address[:len(address)-1], ":"))
			}
		}
	}
	return result
}

func (i *InterfaceMapHandlerImpl) GetServiceUniqueKeyHostAddrMapFromPilot() (map[string]*DubboDetail, error) {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/debug/adsz", i.istioDebugAddr), nil)
	token, err := ioutil.ReadFile(i.istioTokenPath)
	if err != nil {
		logger.Warnf("token of istiod read from %s failed", i.istioTokenPath)
	} else {
		req.Header.Add(authorizationHeader, istiodTokenPrefix+string(token))
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Info("[XDS Wrapped Client] Try getting interface host map from istio IP %s with error %s\n",
			i.istioDebugAddr, err)
		return nil, err
	}

	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	adszRsp := &ADSZResponse{}
	if err := json.Unmarshal(data, adszRsp); err != nil {
		return nil, err
	}
	return adszRsp.GetMap(), nil
}

// NewInterfaceMapHandlerImpl default istioDebugAddr is DefaultIstiodTokenPath
func NewInterfaceMapHandlerImpl(istioDebugAddr, istioTokenPath string) *InterfaceMapHandlerImpl {
	if len(istioTokenPath) == 0 {
		istioTokenPath = DefaultIstiodTokenPath
	}
	return &InterfaceMapHandlerImpl{
		istioDebugAddr: istioDebugAddr,
		istioTokenPath: istioTokenPath,
	}
}
