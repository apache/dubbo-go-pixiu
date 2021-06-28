package registry

type BaseRegistry struct {
	listeners []Listener
}

func NewBaseRegistry() *BaseRegistry {
	return &BaseRegistry{
		listeners: []Listener{},
	}
}
