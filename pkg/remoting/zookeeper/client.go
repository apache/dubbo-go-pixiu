package zookeeper

import (
	"strings"
	"time"
	"sync"
)

import (
	"github.com/pkg/errors"
	"github.com/dubbogo/go-zookeeper/zk"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Options defines the client option.
type Options struct {
	zkName string
	client *ZooKeeperClient
	ts *zk.TestCluster
}

// Option defines the function to load the options
type Option func(*Options)

// WithZkName sets zk client name
func WithZkName(name string) Option {
	return func(opt *Options) {
		opt.zkName = name
	}
}

// ZooKeeperClient represents zookeeper client Configuration
type ZooKeeperClient struct {
	name         string
	ZkAddrs      []string
	sync.RWMutex // for conn
	conn         *zk.Conn
	Timeout      time.Duration
	exit         chan struct{}
	Wait         sync.WaitGroup

	eventRegistry     map[string][]*chan struct{}
	eventRegistryLock sync.RWMutex
}

func NewZooKeeperClient(name string, zkAddrs []string, timeout time.Duration) (*ZooKeeperClient, <-chan zk.Event, error) {
	var (
		err   error
		event <-chan zk.Event
		z     *ZooKeeperClient
	)

	z = &ZooKeeperClient{
		name:          name,
		ZkAddrs:       zkAddrs,
		Timeout:       timeout,
		exit:          make(chan struct{}),
		eventRegistry: make(map[string][]*chan struct{}),
	}
	// connect to zookeeper
	z.conn, event, err = zk.Connect(zkAddrs, timeout)
	if err != nil {
		return nil, nil, errors.WithMessagef(err, "zk.Connect(zkAddrs:%+v)", zkAddrs)
	}

	return z, event, nil
}

// nolint
func StateToString(state zk.State) string {
	switch state {
	case zk.StateDisconnected:
		return "zookeeper disconnected"
	case zk.StateConnecting:
		return "zookeeper connecting"
	case zk.StateAuthFailed:
		return "zookeeper auth failed"
	case zk.StateConnectedReadOnly:
		return "zookeeper connect readonly"
	case zk.StateSaslAuthenticated:
		return "zookeeper sasl authenticated"
	case zk.StateExpired:
		return "zookeeper connection expired"
	case zk.StateConnected:
		return "zookeeper connected"
	case zk.StateHasSession:
		return "zookeeper has session"
	case zk.StateUnknown:
		return "zookeeper unknown state"
	case zk.State(zk.EventNodeDeleted):
		return "zookeeper node deleted"
	case zk.State(zk.EventNodeDataChanged):
		return "zookeeper node data changed"
	default:
		return state.String()
	}
}

func (z *ZooKeeperClient) RegisterHandler(event <-chan zk.Event) {
	z.Wait.Add(1)
	go z.HandleZkEvent(event)
}

// HandleZkEvent handles zookeeper events
func (z *ZooKeeperClient) HandleZkEvent(session <-chan zk.Event) {
	var (
		state int
		event zk.Event
	)

	defer func() {
		z.Wait.Done()
		logger.Infof("zk{path:%v, name:%s} connection goroutine game over.", z.ZkAddrs, z.name)
	}()

	for {
		select {
		case <-z.exit:
			return
		case event = <-session:
			logger.Infof("client{%s} get a zookeeper event{type:%s, server:%s, path:%s, state:%d-%s, err:%v}",
				z.name, event.Type, event.Server, event.Path, event.State, StateToString(event.State), event.Err)
			switch (int)(event.State) {
			case (int)(zk.StateDisconnected):
				logger.Warnf("zk{addr:%s} state is StateDisconnected, so close the zk client{name:%s}.", z.ZkAddrs, z.name)
				z.stop()
				z.Lock()
				conn := z.conn
				z.conn = nil
				z.Unlock()
				if conn != nil {
					conn.Close()
				}
				return
			case (int)(zk.EventNodeDataChanged), (int)(zk.EventNodeChildrenChanged):
				logger.Infof("zkClient{%s} get zk node changed event{path:%s}", z.name, event.Path)
				z.eventRegistryLock.RLock()
				for p, a := range z.eventRegistry {
					if strings.HasPrefix(p, event.Path) {
						logger.Infof("send event{state:zk.EventNodeDataChange, Path:%s} notify event to path{%s} related listener",
							event.Path, p)
						for _, e := range a {
							*e <- struct{}{}
						}
					}
				}
				z.eventRegistryLock.RUnlock()
			case (int)(zk.StateConnecting), (int)(zk.StateConnected), (int)(zk.StateHasSession):
				if state == (int)(zk.StateHasSession) {
					continue
				}
				z.eventRegistryLock.RLock()
				if a, ok := z.eventRegistry[event.Path]; ok && 0 < len(a) {
					for _, e := range a {
						*e <- struct{}{}
					}
				}
				z.eventRegistryLock.RUnlock()
			}
			state = (int)(event.State)
		}
	}
}

// getConn gets zookeeper connection safely
func (z *ZooKeeperClient) getConn() *zk.Conn {
	z.RLock()
	defer z.RUnlock()
	return z.conn
}

func (z *ZooKeeperClient) stop() bool {
	select {
	case <-z.exit:
		return true
	default:
		close(z.exit)
	}

	return false
}


// GetChildren gets children by @path
func (z *ZooKeeperClient) GetChildren(path string) ([]string, error) {
	var (
		children []string
		stat     *zk.Stat
	)
	conn := z.getConn()
	if conn == nil {
		return nil, errors.New("ZooKeeper client has no connection")
	}
	children, stat, err := conn.Children(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, errors.Errorf("path{%s} does not exist", path)
		}
		logger.Errorf("zk.Children(path{%s}) = error(%v)", path, errors.WithStack(err))
		return nil, errors.WithMessagef(err, "zk.Children(path:%s)", path)
	}
	if stat.NumChildren == 0 {
		return nil, errors.Errorf("path{%s} has none children", path)
	}

	return children, nil
}
