package api_load

var ApiLoadTypeMap map[string]ApiLoad

type ApiLoad interface {
	// 第一次初始化加载
	InitLoad() error
	// 后面动态更新加载
	HotLoad() error
	// 清除所有加载的api
	Clear() error
}
