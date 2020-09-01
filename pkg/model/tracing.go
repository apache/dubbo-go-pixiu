package model

// Tracing
type Tracing struct {
	Http Http `yaml:"http" json:"http,omitempty"`
}

// Tracing
type Http struct {
	Name   string      `yaml:"name"`
	Config interface{} `yaml:"config"`
}
