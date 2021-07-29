package etcd

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/coreos/etcd/embed"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

const defaultEtcdV3WorkDir = "/tmp/default-dubbo-go-registry.etcd"

func initEtcd(t *testing.T) *embed.Etcd {
	DefaultListenPeerURLs := "http://localhost:2380"
	DefaultListenClientURLs := "http://localhost:2379"
	lpurl, _ := url.Parse(DefaultListenPeerURLs)
	lcurl, _ := url.Parse(DefaultListenClientURLs)
	cfg := embed.NewConfig()
	cfg.LPUrls = []url.URL{*lpurl}
	cfg.LCUrls = []url.URL{*lcurl}
	cfg.Dir = defaultEtcdV3WorkDir
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TestNewETCDV3Registry(t *testing.T) {
	e := initEtcd(t)

	regConfig := model.Registry{
		Protocol: "etcdv3",
		Timeout:  "2s",
		Address:  "127.0.0.1:2379",
	}
	reg, err := newETCDV3Registry(regConfig)
	assert.Nil(t, err)
	assert.NotNil(t, reg)

	regConfig = model.Registry{
		Protocol: "etcdv3",
		Timeout:  "2xxxxxx",
		Address:  "127.0.0.1:2379",
	}
	reg, err = newETCDV3Registry(regConfig)
	assert.Nil(t, reg)
	assert.NotNil(t, err)

	regConfig = model.Registry{
		Protocol: "etcdv3",
		Timeout:  "2s",
		Address:  "",
	}
	reg, err = newETCDV3Registry(regConfig)
	assert.Nil(t, reg)
	assert.NotNil(t, err)

	e.Close()
}
