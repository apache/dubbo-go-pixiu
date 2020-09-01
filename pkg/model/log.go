package model

// AccessLog
type AccessLog struct {
	Name   string          `yaml:"name" json:"name" mapstructure:"name"`
	Filter AccessLogFilter `yaml:"filter" json:"filter" mapstructure:"filter"`
	Config interface{}     `yaml:"config" json:"config" mapstructure:"config"`
}

// AccessLogFilter
type AccessLogFilter struct {
	StatusCodeFilter StatusCodeFilter `yaml:"status_code_filter" json:"status_code_filter" mapstructure:"status_code_filter"`
	DurationFilter   DurationFilter   `yaml:"duration_filter" json:"duration_filter" mapstructure:"duration_filter"`
}

type StatusCodeFilter struct {
}

type DurationFilter struct {
}
