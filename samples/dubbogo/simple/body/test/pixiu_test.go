package test

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestPost(t *testing.T) {
	var url = "http://localhost:8881/api/v1/test-dubbo/user"
	data := "{\"id\":\"0003\",\"code\":3,\"name\":\"dubbogo\",\"age\":99}"
	contentType := "application/json"
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(url, contentType, bytes.NewBufferString(data))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}
