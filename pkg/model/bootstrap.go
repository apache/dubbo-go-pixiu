package model

// Bootstrap the door
type Bootstrap struct {
	StaticResources  StaticResources  `yaml:"static_resources" json:"static_resources" mapstructure:"static_resources"`
	DynamicResources DynamicResources `yaml:"dynamic_resources" json:"dynamic_resources" mapstructure:"dynamic_resources"`
	Tracing          Tracing          `yaml:"tracing" json:"tracing" mapstructure:"tracing"`
}

// GetListeners
func (bs *Bootstrap) GetListeners() []Listener {
	return bs.StaticResources.Listeners
}

// ExistCluster
func (bs *Bootstrap) ExistCluster(name string) bool {
	if len(bs.StaticResources.Clusters) > 0 {
		for _, v := range bs.StaticResources.Clusters {
			if v.Name == name {
				return true
			}
		}
	}

	return false
}

// StaticResources
type StaticResources struct {
	Listeners []Listener `yaml:"listeners" json:"listeners" mapstructure:"listeners"`
	Clusters  []Cluster  `yaml:"clusters" json:"clusters" mapstructure:"clusters"`
}

// DynamicResources TODO
type DynamicResources struct {
	ApiConfig ApiConfig `yaml:"api_config" json:"apiConfig"`
}
