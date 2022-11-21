/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package zookeeper

import (
	"strings"
	"sync"
	"time"
)

import (
	"github.com/dubbogo/go-zookeeper/zk"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

var (
	// ErrNilZkClientConn no conn error
	ErrNilZkClientConn = errors.New("zookeeper Client{conn} is nil")
	// ErrNilChildren no children error
	ErrNilChildren = errors.Errorf("has none children")
	// ErrNilNode no node error
	ErrNilNode = errors.Errorf("node does not exist")
)

// Options defines the client option.
type Options struct {
	zkName string
	ts     *zk.TestCluster
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

	eventRegistry     map[string][]chan zk.Event
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
		eventRegistry: make(map[string][]chan zk.Event),
	}
	// connect to zookeeper
	z.conn, event, err = zk.Connect(zkAddrs, timeout)
	if err != nil {
		return nil, nil, errors.WithMessagef(err, "zk.Connect(zkAddrs:%+v)", zkAddrs)
	}

	return z, event, nil
}

// StateToString will transfer zk state to string
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
		return "zookeeper has Session"
	case zk.StateUnknown:
		return "zookeeper unknown state"
	default:
		return state.String()
	}
}

func (z *ZooKeeperClient) RegisterHandler(event <-chan zk.Event) {
	z.Wait.Add(1)
	go z.HandleZkEvent(event)
}

// ExistW is a wrapper to the zk connection conn.ExistsW
func (z *ZooKeeperClient) ExistW(zkPath string) (<-chan zk.Event, error) {
	conn := z.getConn()
	if conn == nil {
		return nil, ErrNilZkClientConn
	}
	exist, _, watcher, err := conn.ExistsW(zkPath)
	if err != nil {
		return nil, errors.WithMessagef(err, "zk.ExistsW(path:%s)", zkPath)
	}
	if !exist {
		return nil, errors.WithMessagef(ErrNilNode, "zkClient{%s} App zk path{%s} does not exist.", z.name, zkPath)
	}
	return watcher.EvtCh, nil
}

// HandleZkEvent handles zookeeper events
func (z *ZooKeeperClient) HandleZkEvent(s <-chan zk.Event) {
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
		case event = <-s:
			logger.Infof("client{%s} get a zookeeper event{type:%s, server:%s, path:%s, state:%d-%s, err:%v}",
				z.name, event.Type, event.Server, event.Path, event.State, StateToString(event.State), event.Err)
			switch event.State {
			case zk.StateDisconnected:
				logger.Warnf("zk{addr:%s} state is StateDisconnected, so close the zk client{name:%s}.", z.ZkAddrs, z.name)
				z.Destroy()
				return
			case zk.StateConnected:
				logger.Infof("zkClient{%s} get zk node changed event{path:%s}", z.name, event.Path)
				z.eventRegistryLock.RLock()
				for path, a := range z.eventRegistry {
					if strings.HasPrefix(event.Path, path) {
						logger.Infof("send event{state:zk.EventNodeDataChange, Path:%s} notify event to path{%s} related listener",
							event.Path, path)
						for _, e := range a {
							e <- event
						}
					}
				}
				z.eventRegistryLock.RUnlock()
			case zk.StateConnecting, zk.StateHasSession:
				if state == (int)(zk.StateHasSession) {
					continue
				}
				z.eventRegistryLock.RLock()
				if a, ok := z.eventRegistry[event.Path]; ok && 0 < len(a) {
					for _, e := range a {
						e <- event
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
	conn := z.getConn()
	if conn == nil {
		return nil, ErrNilZkClientConn
	}
	children, stat, err := conn.Children(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, errors.WithMessagef(ErrNilNode, "path{%s} does not exist", path)
		}
		logger.Errorf("zk.Children(path{%s}) = error(%v)", path, errors.WithStack(err))
		return nil, errors.WithMessagef(err, "zk.Children(path:%s)", path)
	}
	if stat.NumChildren == 0 {
		return nil, ErrNilChildren
	}
	return children, nil
}

// GetChildrenW gets children watch by @path
func (z *ZooKeeperClient) GetChildrenW(path string) ([]string, <-chan zk.Event, error) {
	conn := z.getConn()
	if conn == nil {
		return nil, nil, ErrNilZkClientConn
	}
	children, stat, watcher, err := conn.ChildrenW(path)
	if err != nil {
		if err == zk.ErrNoChildrenForEphemerals {
			return nil, nil, ErrNilChildren
		}
		if err == zk.ErrNoNode {
			return nil, nil, ErrNilNode
		}
		return nil, nil, errors.WithMessagef(err, "zk.ChildrenW(path:%s)", path)
	}
	if stat == nil {
		return nil, nil, errors.Errorf("path{%s} get stat is nil", path)
	}
	//if len(children) == 0 {
	//	return nil, nil, ErrNilChildren
	//}

	return children, watcher.EvtCh, nil
}

// GetContent gets content by @path
func (z *ZooKeeperClient) GetContent(path string) ([]byte, error) {
	var (
		data []byte
	)
	conn := z.getConn()
	if conn == nil {
		return nil, errors.New("ZooKeeper client has no connection")
	}
	data, _, err := conn.Get(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, errors.Errorf("path{%s} does not exist", path)
		}
		logger.Errorf("zk.Data(path{%s}) = error(%v)", path, errors.WithStack(err))
		return nil, errors.WithMessagef(err, "zk.Data(path:%s)", path)
	}
	return data, nil
}

func (z *ZooKeeperClient) GetConnState() zk.State {
	conn := z.getConn()
	if conn != nil {
		return conn.State()
	}
	return zk.StateExpired
}

func (z *ZooKeeperClient) Destroy() {
	z.stop()
	z.Lock()
	conn := z.conn
	z.conn = nil
	z.Unlock()
	if conn != nil {
		conn.Close()
	}
}
