package extension

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

type (
	AdapterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string

		// CreateAdapter return the Adapter callback
		CreateAdapter(server *server.Server, config interface{}, bs *model.Bootstrap) (Adapter, error)
	}

	Adapter interface {
	}
)
