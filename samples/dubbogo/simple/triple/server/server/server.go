package main

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"dubbo.apache.org/dubbo-go/v3/config"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
)

import (
	"github.com/apache/dubbo-go-pixiu/samples/dubbogo/simple/triple/server/api"
)

func main() {
	config.SetProviderService(&GreeterProvider{})

	rc := config.NewRootConfigBuilder().
		SetProvider(config.NewProviderConfigBuilder().
			AddService("GreeterProvider", config.NewServiceConfigBuilder().Build()).
			Build()).
		AddProtocol("tripleProtocolKey", config.NewProtocolConfigBuilder().
			SetName("tri").
			SetPort("20001").
			Build()).
		AddRegistry("registryKey", config.NewRegistryConfigWithProtocolDefaultPort("zookeeper")).
		Build()

	// start dubbo-go framework with configuration
	if err := config.Load(config.WithRootConfig(rc)); err != nil{
		panic(err)
	}

	select {}
}

type GreeterProvider struct {
	api.UnimplementedUserProviderServer
}

func (s *GreeterProvider) SayHelloStream(svr api.UserProvider_SayHelloStreamServer) error {
	c, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go GreeterProvider recv 1 user, name = %s\n", c.Name)
	c2, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go GreeterProvider recv 2 user, name = %s\n", c2.Name)
	c3, err := svr.Recv()
	if err != nil {
		return err
	}
	logger.Infof("Dubbo-go GreeterProvider recv 3 user, name = %s\n", c3.Name)

	if err:= svr.Send(&api.User{
		Name: "hello " + c.Name,
		Age:  18,
		Id:   "123456789",
	}); err != nil{
		return err
	}
	if err := svr.Send(&api.User{
		Name: "hello " + c2.Name,
		Age:  19,
		Id:   "123456789",
	}); err != nil{
		return err
	}
	return nil
}

func (s *GreeterProvider) SayHello(ctx context.Context, in *api.HelloRequest) (*api.User, error) {
	logger.Infof("Dubbo3 GreeterProvider get user name = %s\n", in.Name)
	return &api.User{Name: "Hello " + in.Name, Id: "12345", Age: 21}, nil
}