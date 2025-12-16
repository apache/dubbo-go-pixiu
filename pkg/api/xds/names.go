package xds

const (
	ResourceTypePrefix = "dubbo-go.pixiu"
	ClusterType        = ResourceTypePrefix + "/v1/discovery:cluster"
	ListenerType       = ResourceTypePrefix + "/v1/discovery:listener"
	EndpointType       = ResourceTypePrefix + "/v1/discovery:endpoint"
	RouterType         = ResourceTypePrefix + "/v1/discovery:route"
	RuntimeType        = ResourceTypePrefix + "/v1/discovery:runtime"
)
