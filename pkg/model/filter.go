package model

// FilterChain filter chain
type FilterChain struct {
	FilterChainMatch FilterChainMatch `yaml:"filter_chain_match" json:"filter_chain_match" mapstructure:"filter_chain_match"`
	Filters          []Filter         `yaml:"filters" json:"filters" mapstructure:"filters"`
}

// Filter core struct, filter is extend by user
type Filter struct {
	Name   string      `yaml:"name" json:"name" mapstructure:"name"`       // Name filter name unique
	Config interface{} `yaml:"config" json:"config" mapstructure:"config"` // Config filter config
}

// FilterChainMatch
type FilterChainMatch struct {
	Domains []string `yaml:"domains" json:"domains" mapstructure:"domains"`
}
