package model

// HealthCheck
type HealthCheck struct {
}

// HttpHealthCheck
type HttpHealthCheck struct {
	HealthCheck
	Host             string
	Path             string
	UseHttp2         bool
	ExpectedStatuses int64
}

// GrpcHealthCheck
type GrpcHealthCheck struct {
	HealthCheck
	ServiceName string
	Authority   string
}

// CustomHealthCheck
type CustomHealthCheck struct {
	HealthCheck
	Name   string
	Config interface{}
}
