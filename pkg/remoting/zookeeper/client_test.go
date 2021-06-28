package zookeeper

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/pkg/errors"
	"github.com/dubbogo/go-zookeeper/zk"
	"testing"
)

func verifyEventStateOrder(t *testing.T, c <-chan zk.Event, expectedStates []zk.State, source string) {
	for _, state := range expectedStates {
		for {
			event, ok := <-c
			if !ok {
				t.Fatalf("unexpected channel close for %s", source)
			}
			logger.Debug(event)
			if event.Type != zk.EventSession {
				continue
			}
			assert.Equal(t, event.State, state)
			break
		}
	}
}

// NewMockZooKeeperClient returns a mock client instance
func NewMockZooKeeperClient(name string, timeout time.Duration, opts ...Option) (*zk.TestCluster, *ZooKeeperClient, <-chan zk.Event, error) {
	var (
		err   error
		event <-chan zk.Event
		z     *ZooKeeperClient
		ts    *zk.TestCluster
	)

	z = &ZooKeeperClient{
		name:          name,
		ZkAddrs:       []string{},
		Timeout:       timeout,
		exit:          make(chan struct{}),
		eventRegistry: make(map[string][]*chan struct{}),
	}

	options := &Options{}
	for _, opt := range opts {
		opt(options)
	}

	// connect to zookeeper
	if options.ts != nil {
		ts = options.ts
	} else {
		ts, err = zk.StartTestCluster(1, nil, nil)
		if err != nil {
			return nil, nil, nil, errors.WithMessagef(err, "zk.Connect")
		}
	}

	z.conn, event, err = ts.ConnectWithOptions(timeout)
	if err != nil {
		return nil, nil, nil, errors.WithMessagef(err, "zk.Connect")
	}

	return ts, z, event, nil
}

func TestNewZooKeeperClient(t *testing.T) {
	tc, err := zk.StartTestCluster(1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tc.Stop()

	hosts := make([]string, len(tc.Servers))
	for i, srv := range tc.Servers {
		hosts[i] = fmt.Sprintf("127.0.0.1:%d", srv.Port)
	}
	zkClient, eventChan, err := NewZooKeeperClient("testZK", hosts, 30 * time.Second)
	if err != nil {
		t.Fatalf("Connect returned error: %+v", err)
	}

	states := []zk.State{zk.StateConnecting, zk.StateConnected, zk.StateHasSession}
	verifyEventStateOrder(t, eventChan, states, "event channel")

	zkClient.getConn().Close()
	verifyEventStateOrder(t, eventChan, []zk.State{zk.StateDisconnected}, "event channel")
}

func TestGetChildren(t *testing.T) {
	tc, err := zk.StartTestCluster(1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tc.Stop()

	conn, _, err := tc.ConnectAll()
	assert.Nil(t, err)
	path, err := conn.Create("/test", nil, 0, zk.WorldACL(zk.PermAll))
	assert.Nil(t, err)
	assert.NotNil(t, path)
	path, err = conn.Create("/test/testchild1", nil, 0, zk.WorldACL(zk.PermAll))
	assert.Nil(t, err)
	assert.NotNil(t, path)
	conn.Create("/test/testchild2", nil, 0, zk.WorldACL(zk.PermAll))

	hosts := make([]string, len(tc.Servers))
	for i, srv := range tc.Servers {
		hosts[i] = fmt.Sprintf("127.0.0.1:%d", srv.Port)
	}
	zkClient, eventChan, err := NewZooKeeperClient("testZK", hosts, 30 * time.Second)
	assert.Nil(t, err)
	wait:
		for {
			event := <- eventChan
			switch event.State {
			case zk.StateDisconnected:
				break wait
			case zk.StateConnected:
				children, err := zkClient.GetChildren("/test")
				assert.Nil(t, err)
				assert.Equal(t, children[1], "testchild1")
				assert.Equal(t, children[0], "testchild2")

				children, err = zkClient.GetChildren("/vacancy")
				assert.EqualError(t, err, "path{/vacancy} does not exist")
				assert.Nil(t, children)
				children, err = zkClient.GetChildren("/test/testchild1")
				assert.EqualError(t, err, "path{/test/testchild1} has none children")
				assert.Empty(t, children)
				zkClient.conn.Close()
			}
		}
}
