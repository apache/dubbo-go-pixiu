package service

// DiscoveryRequest a request for discovery
type DiscoveryRequest struct {
	Body []byte
}

// NewDiscoveryRequest return a DiscoveryRequest with body
func NewDiscoveryRequest(b []byte) *DiscoveryRequest {
	return &DiscoveryRequest{
		Body: b,
	}
}

// DiscoveryResponse a response for discovery
type DiscoveryResponse struct {
	Success bool
	Data    interface{}
}

// NewDiscoveryResponseWithSuccess return a DiscoveryResponse with success
func NewDiscoveryResponseWithSuccess(b bool) *DiscoveryResponse {
	return &DiscoveryResponse{
		Success: b,
	}
}

// NewDiscoveryResponse return a DiscoveryResponse with Data and success true
func NewDiscoveryResponse(d interface{}) *DiscoveryResponse {
	return &DiscoveryResponse{
		Success: true,
		Data:    d,
	}
}

var EmptyDiscoveryResponse = &DiscoveryResponse{}

// ApiDiscoveryService api discovery service interface
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
