package discovery

// 封装服务实例信息
type InstanceInfo struct {

	ServiceName		string
	// 格式为 host:port
	Addr 			string
	// 此实例附加信息
	Meta 			map[string]string
}

// 服务发现客户端接口
type Client interface {
	// 直接向远程注册中心查询所有服务实例
	QueryServices() ([]*InstanceInfo, error)

	// 注册自己
	Register() error

	// 取消注册自己
	UnRegister() error

	// 从本地缓存中查询指定服务的全部实例信息
	Get(string) []*InstanceInfo

	// 启动注册信息定时刷新逻辑
	StartPeriodicalRefresh() error

	// 获取内部保存的注册表
	GetInternalRegistryStore() *InsInfoArrSyncMap

	// 更新内部保存的注册表
	SetInternalRegistryStore(*InsInfoArrSyncMap)
}