package server

import "github.com/apache/dubbo-go-pixiu/pkg/model"

type AdapterManager struct {
	cConfig  []*model.ClusterAdapter
	caConfig []*model.ClusterAdapter
}
