package service

type DiscoveryRequest struct {
	Body []byte
}

func NewDiscoveryRequest(b []byte) *DiscoveryRequest {
	return &DiscoveryRequest{
		Body: b,
	}
}

type DiscoveryResponse struct {
	Success bool
	Data    interface{}
}

func NewSuccessDiscoveryResponse() *DiscoveryResponse {
	return &DiscoveryResponse{
		Success: true,
	}
}

func NewDiscoveryResponse(d interface{}) *DiscoveryResponse {
	return &DiscoveryResponse{
		Success: true,
		Data:    d,
	}
}

var EmptyDiscoveryResponse = &DiscoveryResponse{}

type ApiDiscoveryService interface {
	AddApi(request DiscoveryRequest) (DiscoveryResponse, error)
	GetApi(request DiscoveryRequest) (DiscoveryResponse, error)
}

type ListenerDiscoveryService interface {
	AddListeners(request DiscoveryRequest) (DiscoveryResponse, error)
	GetListeners(request DiscoveryRequest) (DiscoveryResponse, error)
}

type RouteDiscoveryService interface {
	AddRoutes(r DiscoveryRequest) (DiscoveryResponse, error)
	GetRoutes(r DiscoveryRequest) (DiscoveryResponse, error)
}

type ClusterDiscoveryService interface {
	AddClusters(r DiscoveryRequest) (DiscoveryResponse, error)
	GetClusters(r DiscoveryRequest) (DiscoveryResponse, error)
}

type EndpointDiscoveryService interface {
	AddEndpoints(r DiscoveryRequest) (DiscoveryResponse, error)
	GetEndpoints(r DiscoveryRequest) (DiscoveryResponse, error)
}
