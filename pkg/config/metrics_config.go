package config

// Metrics provides options to expose and send dubbo-go-proxy metrics for Prometheus
// Other third party monitoring systems is not support.
type Metrics struct {
	Prometheus *Prometheus `yaml:"prometheus"`
}

// Prometheus can contain specific configuration used by the Prometheus Metrics exporter.
type Prometheus struct {
	Buckets              []float64 `description:"Buckets for latency metrics." json:"buckets,omitempty" toml:"buckets,omitempty" yaml:"buckets,omitempty" export:"true"`
	AddEntryPointsLabels bool      `description:"Enable metrics on entry points." json:"addEntryPointsLabels,omitempty" toml:"addEntryPointsLabels,omitempty" yaml:"addEntryPointsLabels,omitempty" export:"true"`
	AddServicesLabels    bool      `description:"Enable metrics on services." json:"addServicesLabels,omitempty" toml:"addServicesLabels,omitempty" yaml:"addServicesLabels,omitempty" export:"true"`
	EntryPoint           string    `description:"EntryPoint" export:"true" json:"entryPoint,omitempty" toml:"entryPoint,omitempty" yaml:"entryPoint,omitempty"`
}


