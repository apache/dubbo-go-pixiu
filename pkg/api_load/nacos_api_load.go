package api_load

func init() {
	var _ ApiLoad = new(NacosApiLoader)
}

type NacosApiLoader struct {
}

func (f *NacosApiLoader) InitLoad() (err error) {
	panic("")
}

func (f *NacosApiLoader) HotLoad() error {
	panic("")
}

func (f *NacosApiLoader) Clear() error {
	panic("")
}
