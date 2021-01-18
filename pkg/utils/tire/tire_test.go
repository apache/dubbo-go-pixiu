package tire

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	tire := NewTire()
	ret := tire.Put("/path1/{pathvarible1}/path2/{pathvarible2}", nil)
	assert.True(t, ret)
	ret = tire.Put("/path2/{pathvarible1}/path2/{pathvarible2}", nil)
	assert.True(t, ret)
	ret = tire.Put("/path2/3/path2/{pathvarible2}", nil)
	assert.True(t, ret)

	ret = tire.Put("/path2/3/path2/{pathvarible2}", nil)
	assert.False(t, ret)

	ret = tire.Put("/path2/3/path2/{pathvarible2}/3", nil)
	assert.True(t, ret)
	ret = tire.Put("/path2/3/path2/{432423}/3", nil)
	assert.False(t, ret)
	ret = tire.Put("/path2/3/path2/{432423}/3/a/b/c/d/{fdsa}", nil)
	assert.True(t, ret)

	ret = tire.Put("/path2/3/path2/{432423}/3/a/b/c/c/{fdsa}", nil)
	assert.True(t, ret)

	ret = tire.Put("/path2/3/path2/{432423}/3/a/b/c/c/{fdsafdsafsdafsda}", nil)
	assert.False(t, ret)

	ret = tire.Put("/path1/{pathvarible1}/path2/{pathvarible2}/{fdsa}", nil)
	assert.True(t, ret)

	ret = tire.Put("/path1/{432}/path2/{34}", nil)
	assert.False(t, ret)
}
