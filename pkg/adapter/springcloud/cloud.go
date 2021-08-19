package springcloud

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.SpringCloudAdapter
)

func init() {
	extension.RegisterAdapterPlugin(&CloudPlugin{})
}

type (
	CloudPlugin struct {
	}

	CloudAdapter struct {
		cfg    *Config
		server *server.Server
	}

	Config struct {
		Registry *Registry `yaml:"registry" json:"registry" default:"registry"`
	}

	// Registry remote registry where dubbo apis are registered.
	Registry struct {
		Protocol string `yaml:"protocol" json:"protocol" default:"zookeeper"`
		Timeout  string `yaml:"timeout" json:"timeout"`
		Address  string `yaml:"address" json:"address"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`
	}
)

func (p *CloudPlugin) Kind() string {
	return Kind
}

func (p *CloudPlugin) CreateAdapter(server *server.Server, config interface{}, bs *model.Bootstrap) (extension.Adapter, error) {
	specConfig := config.(*Config)
	return &CloudAdapter{cfg: specConfig, server: server}, nil
}

func (a *CloudAdapter) Start() {
	// do not block the main goroutine
	go func() {
		// fetch service instance from consul

		// transform into endpoint and cluster
		endpoint := &model.Endpoint{}
		endpoint.ID = "user-mock-service"
		endpoint.Address = model.SocketAddress{
			Address: "127.0.0.1",
			Port:    8080,
		}
		cluster := &model.Cluster{}
		cluster.Name = "userservice"
		cluster.Lb = model.Rand
		cluster.Endpoints = []*model.Endpoint{
			endpoint,
		}
		// add cluster into manager
		cm := a.server.GetClusterManager()
		cm.AddCluster(cluster)

		// transform into route
		routeMatch := model.RouterMatch{
			Prefix: "/user",
		}
		routeAction := model.RouteAction{
			Cluster: "userservice",
		}
		route := &model.Router{Match: routeMatch, Route: routeAction}

		a.server.GetRouterManager().AddRouter(route)

	}()
}

func (a *CloudAdapter) Stop() {

}
