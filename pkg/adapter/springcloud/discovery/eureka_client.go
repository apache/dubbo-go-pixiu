package discovery

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"net"
	"strconv"
	"time"

	"github.com/wanghongfei/go-eureka-client/eureka"
)

// var euClient *eureka.Client
var gogateApp *eureka.InstanceInfo
var instanceId = ""

var ticker *time.Ticker
var tickerCloseChan chan struct{}

type EurekaClient struct {
	// 继承方法
	*periodicalRefreshClient

	client *eureka.Client

	// 保存服务地址
	// key: 服务名:版本号, 版本号为eureka注册信息中的metadata[version]值
	// val: []*InstanceInfo
	registryMap *InsInfoArrSyncMap
}

func NewEurekaClient(confFile string) (DiscoveryClient, error) {
	c, err := eureka.NewClientFromFile(confFile)
	if nil != err {
		logger.Errorf("failed to init eureka client", err)
		return nil, err
	}

	euClient := &EurekaClient{client: c}
	euClient.periodicalRefreshClient = newPeriodicalRefresh(euClient)

	return euClient, nil
}

func (c *EurekaClient) Get(serviceId string) []*InstanceInfo {
	instance, exist := c.registryMap.Get(serviceId)
	if !exist {
		return nil
	}

	return instance
}

func (c *EurekaClient) GetInternalRegistryStore() *InsInfoArrSyncMap {
	return c.registryMap
}

func (c *EurekaClient) SetInternalRegistryStore(registry *InsInfoArrSyncMap) {
	c.registryMap = registry
}

// QueryServices 查询所有服务
func (c *EurekaClient) QueryServices() ([]*InstanceInfo, error) {
	apps, err := c.client.GetApplications()
	if nil != err {
		logger.Errorf("faield to query eureka", err)
		return nil, err
	}

	var instances []*InstanceInfo
	for _, app := range apps.Applications {
		// 服务名
		servName := app.Name

		// 遍历每一个实例
		for _, ins := range app.Instances {
			// 跳过无效实例
			if nil == ins.Port || ins.Status != "UP" {
				continue
			}

			addr := ins.HostName + ":" + strconv.Itoa(ins.Port.Port)
			var meta map[string]string
			if nil != ins.Metadata {
				meta = ins.Metadata.Map
			}

			instances = append(
				instances,
				&InstanceInfo{
					ServiceName: servName,
					Addr:        addr,
					Meta:        meta,
				},
			)
		}
	}

	return instances, nil
}

func (c *EurekaClient) Register() error {
	ip, err := GetFirstNoneLoopIp()
	if nil != err {
		logger.Errorf("failed to get first none loop ip", err)
		return err
	}

	instanceId = ip + ":" + strconv.Itoa(9909)

	// 注册
	logger.Infof("register to eureka as %s", instanceId)
	gogateApp = eureka.NewInstanceInfo(
		instanceId,
		"Pixiu",
		ip,
		9909,
		30,
		false,
	)
	gogateApp.Metadata = &eureka.MetaData{
		Class: "",
		Map:   map[string]string{"version": "1.0.0"},
	}

	err = c.client.RegisterInstance("Pixiu", gogateApp)
	if nil != err {
		logger.Errorf("failed to register to eureka", err)
		return err
	}

	// 心跳
	go func() {
		ticker = time.NewTicker(time.Second * time.Duration(30))
		tickerCloseChan = make(chan struct{})

		for {
			select {
			case <-ticker.C:
				c.heartbeat()

			case <-tickerCloseChan:
				logger.Info("heartbeat stopped")
				return

			}
		}
	}()

	return nil
}

func (c *EurekaClient) UnRegister() error {
	c.stopHeartbeat()

	logger.Infof("unregistering %s", instanceId)
	err := c.client.UnregisterInstance("Pixiu", instanceId)

	if nil != err {
		logger.Errorf("failed to unregister", err)
		return err
	}

	logger.Info("done unregistration")
	return nil
}

func (c *EurekaClient) stopHeartbeat() {
	ticker.Stop()
	close(tickerCloseChan)
}

func (c *EurekaClient) heartbeat() {
	err := c.client.SendHeartbeat(gogateApp.App, instanceId)
	if nil != err {
		logger.Warnf("failed to send heartbeat, %v", err)
		return
	}

	logger.Info("heartbeat sent")
}

func GetFirstNoneLoopIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if nil != err {
		return "", fmt.Errorf("failed to fetch interfaces => %w", err)
	}

	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no first-none-loop ip found")
}
