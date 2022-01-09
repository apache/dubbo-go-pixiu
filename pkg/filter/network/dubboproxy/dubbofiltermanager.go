package dubboproxy

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"sync"
)

// DubboFilterManager manage filters
type DubboFilterManager struct {
	filterConfigs []*model.DubboFilter
	filters       []filter.DubboFilter
	mu            sync.RWMutex
}

// NewDubboFilterManager create filter manager
func NewDubboFilterManager(fs []*model.DubboFilter) *DubboFilterManager {
	filters := createDubboFilter(fs)
	fm := &DubboFilterManager{filterConfigs: fs, filters: filters}
	return fm
}

func createDubboFilter(fs []*model.DubboFilter) []filter.DubboFilter {
	var filters []filter.DubboFilter

	for _, f := range fs {
		p, err := filter.GetDubboFilterPlugin(f.Name)
		if err != nil {
			logger.Error("createDubboFilter %s getNetworkFilterPlugin error %s", f.Name, err)
		}

		config := p.Config()
		if err := yaml.ParseConfig(config, f.Config); err != nil {
			logger.Error("createDubboFilter %s parse config error %s", f.Name, err)
		}

		filter, err := p.CreateFilter(config)
		if err != nil {
			logger.Error("createDubboFilter %s createFilter error %s", f.Name, err)
		}
		filters = append(filters, filter)
	}
	return filters
}
