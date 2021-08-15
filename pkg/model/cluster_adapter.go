package model

type ClusterAdapter struct {
	Name             string              `yaml:"name" json:"name"`             // Name the cluster unique name
	Registries       map[string]Registry `yaml:"registries" json:"registries"`
}
