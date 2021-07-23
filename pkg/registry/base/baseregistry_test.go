package baseregistry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"math/rand"
	"testing"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
	"github.com/stretchr/testify/assert"
)

var _ registry.Registry = new(mockRegFacade)

type mockRegFacade struct {
	BaseRegistry
}

func (m *mockRegFacade) LoadInterfaces() ([]registry.RegisteredAPI, []error) {
	return nil, nil
}
func (m *mockRegFacade) LoadApplications() ([]registry.RegisteredAPI, []error) {
	return nil, nil
}
func (m *mockRegFacade) DoSubscribe(registry.RegisteredAPI) error {
	return nil
}
func (m *mockRegFacade) DoUnsubscribe(registry.RegisteredAPI) error {
	return nil
}
func (m *mockRegFacade) DoSDSubscribe(registry.RegisteredAPI) error {
	return nil
}
func (m *mockRegFacade) DoSDUnsubscribe(registry.RegisteredAPI) error {
	return nil
}

func CreateMockRegisteredAPI(urlPattern string) router.API {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	randStringRunes := func(n int) string {
		rand.Seed(time.Now().UnixNano())
		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		return string(b)
	}
	if len(urlPattern) == 0 {
		urlPattern = randStringRunes(4) + "/" + randStringRunes(4) + "/" + randStringRunes(4)
	}
	return router.API{
		URLPattern: urlPattern,
	}
}

func TestDeduplicate(t *testing.T) {
	var apis []registry.RegisteredAPI
	for i := 0; i < 10; i++ {
		apis = append(apis, registry.RegisteredAPI{
			API: CreateMockRegisteredAPI(""),
		})
	}
	var facade *mockRegFacade = &mockRegFacade{}
	baseReg := NewBaseRegistry(facade)
	rst := baseReg.deduplication(apis)
	assert.Equal(t, len(rst), 10)

	apis = []registry.RegisteredAPI{
		registry.RegisteredAPI{
			API: CreateMockRegisteredAPI("/happy/friday"),
		},
		registry.RegisteredAPI{
			API: CreateMockRegisteredAPI("/happy/friday"),
		},
		registry.RegisteredAPI{
			API: CreateMockRegisteredAPI("/happy/Saturday"),
		},
	}
	rst = baseReg.deduplication(apis)
	assert.Equal(t, len(rst), 2)
}
