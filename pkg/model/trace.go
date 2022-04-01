package model

type TracerConfig struct {
	Name    string                 `yaml:"name" json:"name" mapstructure:"name"`
	Sampler Sampler                `yaml:"sampler" json:"sampler" mapstructure:"sampler"`
	Config  map[string]interface{} `yaml:"config" json:"config" mapstructure:"config"`
}
type Sampler struct {
	Type  string  `yaml:"type" json:"type" mapstructure:"type"`
	Param float64 `yaml:"param" json:"param" mapstructure:"param"`
}
