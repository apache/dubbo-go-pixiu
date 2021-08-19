package model

type Adapter struct {
	Name   string      `yaml:"name" json:"name"`                           // Name the adapter unique name
	Config interface{} `yaml:"config" json:"config" mapstructure:"config"` // Config adapter config
}
