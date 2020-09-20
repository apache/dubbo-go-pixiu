package model

type ApiConfig struct {
	Nacos *Nacos `yaml:"nacos" json:"nacos"`
	File  *File  `yaml:"file" json:"file"`
}

type Nacos struct {
	Address string `yaml:"address" json:"address"`
	Prior   int    `yaml:"prior" json:"prior"`
}

type File struct {
	FileApiConfPath string `yaml:"file-api-conf-path" json:"fileApiConfPath"`
}
