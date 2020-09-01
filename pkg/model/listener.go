package model

// Listener is a server, listener a port
type Listener struct {
	Name    string  `yaml:"name" json:"name" mapstructure:"name"`
	Address Address `yaml:"address" json:"address" mapstructure:"address"`

	FilterChains []FilterChain `yaml:"filter_chains" json:"filter_chains" mapstructure:"filter_chains"`
	Config       interface{}   `yaml:"config" json:"config" mapstructure:"config"`
}
