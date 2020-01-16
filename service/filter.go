package service

var (
	filters = make(map[string]func() Filter)
)

type Filter interface {
	OnRequest(ctx ProxyContext) (ret int)
}

func SetFilter(name string, v func() Filter) {
	filters[name] = v
}

func GetFilter(name string) Filter {
	if filters[name] == nil {
		panic("filter for " + name + " is not existing!")
	}
	return filters[name]()
}
