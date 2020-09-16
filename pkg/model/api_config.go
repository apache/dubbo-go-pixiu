package model

type ApiConfig struct {
	Nacos *Nacos `yaml:"nacos" json:"nacos"`
}

type Nacos struct {
	Address string `yaml:"address" json:"address"`
}
