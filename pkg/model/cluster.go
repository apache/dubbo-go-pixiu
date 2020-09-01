package model

// Cluster a single upstream cluster
type Cluster struct {
	Name              string           `yaml:"name" json:"name"`           // Name the cluster unique name
	TypeStr           string           `yaml:"type" json:"type"`           // Type the cluster discovery type string value
	Type              DiscoveryType    `yaml:"omitempty" json:"omitempty"` // Type the cluster discovery type
	EdsClusterConfig  EdsClusterConfig `yaml:"eds_cluster_config" json:"eds_cluster_config"`
	LbStr             string           `yaml:"lb_policy" json:"lb_policy"`             // Lb the cluster select node used loadBalance policy
	Lb                LbPolicy         `yaml:"omitempty" json:"omitempty"`             // Lb the cluster select node used loadBalance policy
	ConnectTimeoutStr string           `yaml:"connect_timeout" json:"connect_timeout"` // ConnectTimeout timeout for connect to cluster node
	HealthChecks      []HealthCheck    `yaml:"health_checks" json:"health_checks"`
	Hosts             []Address        `yaml:"hosts" json:"hosts"` // Hosts whe discovery type is Static, StrictDNS or LogicalDns，this need config
}

// DiscoveryType
type DiscoveryType int32

const (
	Static DiscoveryType = 0 + iota
	StrictDNS
	LogicalDns
	EDS
	OriginalDst
)

// DiscoveryTypeName
var DiscoveryTypeName = map[int32]string{
	0: "Static",
	1: "StrictDNS",
	2: "LogicalDns",
	3: "EDS",
	4: "OriginalDst",
}

// DiscoveryTypeValue
var DiscoveryTypeValue = map[string]int32{
	"Static":      0,
	"StrictDNS":   1,
	"LogicalDns":  2,
	"EDS":         3,
	"OriginalDst": 4,
}

// EdsClusterConfig
type EdsClusterConfig struct {
	EdsConfig   ConfigSource `yaml:"eds_config" json:"eds_config" mapstructure:"eds_config"`
	ServiceName string       `yaml:"service_name" json:"service_name" mapstructure:"service_name"`
}
