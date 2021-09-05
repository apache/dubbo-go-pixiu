package baseregistry

import (
	"math/rand"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

var _ registry.Registry = new(mockRegFacade)

type mockRegFacade struct {
	BaseRegistry
}

func (m *mockRegFacade) DoSubscribe() error {
	return nil
}
func (m *mockRegFacade) DoUnsubscribe() error {
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
