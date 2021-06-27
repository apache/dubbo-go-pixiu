package zookeeper

import (
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

// NewMockZookeeperClient returns a mock client instance
func NewMockZookeeperClient(name string, timeout time.Duration, opts ...Option) (*zk.TestCluster, *ZookeeperClient, <-chan zk.Event, error) {
	var (
		err   error
		event <-chan zk.Event
		z     *ZookeeperClient
		ts    *zk.TestCluster
	)

	z = &ZookeeperClient{
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

	z.Conn, event, err = ts.ConnectWithOptions(timeout)
	if err != nil {
		return nil, nil, nil, errors.WithMessagef(err, "zk.Connect")
	}

	return ts, z, event, nil
}

func TestNewZookeeperClient(t *testing.T) {
	tc, err := zk.StartTestCluster(1, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tc.Stop()

	callbackChan := make(chan zk.Event)
	f := func(event zk.Event) {
		callbackChan <- event
	}

	zook, eventChan, err := tc.ConnectWithOptions(15*time.Second, zk.WithEventCallback(f))
	if err != nil {
		t.Fatalf("Connect returned error: %+v", err)
	}

	states := []zk.State{zk.StateConnecting, zk.StateConnected, zk.StateHasSession}
	verifyEventStateOrder(t, callbackChan, states, "callback")
	verifyEventStateOrder(t, eventChan, states, "event channel")

	zook.Close()
	verifyEventStateOrder(t, callbackChan, []zk.State{zk.StateDisconnected}, "callback")
	verifyEventStateOrder(t, eventChan, []zk.State{zk.StateDisconnected}, "event channel")

}