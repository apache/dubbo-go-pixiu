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

package healthcheck

import (
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

import (
	gxtime "github.com/dubbogo/gost/time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

const (
	DefaultTimeout  = time.Second
	DefaultInterval = 3 * time.Second

	DefaultHealthyThreshold   uint32 = 5
	DefaultUnhealthyThreshold uint32 = 5
	DefaultFirstInterval             = 5 * time.Second
)

type HealthChecker struct {
	checkers      map[string]*EndpointChecker
	sessionConfig map[string]interface{}
	// check config
	timeout            time.Duration
	intervalBase       time.Duration
	healthyThreshold   uint32
	initialDelay       time.Duration
	cluster            *model.ClusterConfig
	unhealthyThreshold uint32
}

// EndpointChecker is a wrapper of types.HealthCheckSession for health check
type EndpointChecker struct {
	endpoint      *model.Endpoint
	HealthChecker *HealthChecker
	// TCP checker, can extend to http, grpc, dubbo or other protocol checker
	tcpChecker    *TCPChecker
	resp          chan checkResponse
	timeout       chan bool
	checkID       uint64
	stop          chan struct{}
	checkTimer    *gxtime.Timer
	checkTimeout  *gxtime.Timer
	unHealthCount uint32
	healthCount   uint32
	threshold     uint32

	once sync.Once
}

type checkResponse struct {
	ID      uint64
	Healthy bool
}

func CreateHealthCheck(cluster *model.ClusterConfig, cfg model.HealthCheckConfig) *HealthChecker {

	timeout, err := time.ParseDuration(cfg.TimeoutConfig)
	if err != nil {
		logger.Infof("[health check] timeout parse duration error %s", err)
		timeout = DefaultTimeout
	}

	interval, err := time.ParseDuration(cfg.IntervalConfig)
	if err != nil {
		logger.Infof("[health check] interval parse duration error %s", err)
		interval = DefaultInterval
	}

	initialDelay, err := time.ParseDuration(cfg.IntervalConfig)
	if err != nil {
		logger.Infof("[health check] initialDelay parse duration error %s", err)
		initialDelay = DefaultFirstInterval
	}

	unhealthyThreshold := cfg.UnhealthyThreshold
	if unhealthyThreshold == 0 {
		unhealthyThreshold = DefaultUnhealthyThreshold
	}
	healthyThreshold := cfg.HealthyThreshold
	if healthyThreshold == 0 {
		healthyThreshold = DefaultHealthyThreshold
	}

	hc := &HealthChecker{
		sessionConfig:      cfg.SessionConfig,
		cluster:            cluster,
		timeout:            timeout,
		intervalBase:       interval,
		healthyThreshold:   healthyThreshold,
		unhealthyThreshold: unhealthyThreshold,
		initialDelay:       initialDelay,
		checkers:           make(map[string]*EndpointChecker),
	}

	return hc
}

func (hc *HealthChecker) Start() {
	// each endpoint
	for _, h := range hc.cluster.Endpoints {
		hc.startCheck(h)
	}
}

func (hc *HealthChecker) Stop() {
	for _, h := range hc.cluster.Endpoints {
		hc.stopCheck(h)
	}
}

func (hc *HealthChecker) StopOne(endpoint *model.Endpoint) {
	hc.stopCheck(endpoint)
}

func (hc *HealthChecker) StartOne(endpoint *model.Endpoint) {
	hc.startCheck(endpoint)
}

func (hc *HealthChecker) startCheck(endpoint *model.Endpoint) {
	addr := endpoint.Address.GetAddress()
	if _, ok := hc.checkers[addr]; !ok {
		c := newChecker(endpoint, hc)
		hc.checkers[addr] = c
		go c.Start()
		logger.Infof("[health check] create a health check session for %s", addr)
	}
}

func (hc *HealthChecker) stopCheck(endpoint *model.Endpoint) {
	addr := endpoint.Address.GetAddress()
	if c, ok := hc.checkers[addr]; ok {
		c.Stop()
		delete(hc.checkers, addr)
		logger.Infof("[health check] create a health check session for %s", addr)
	}
}

func newChecker(endpoint *model.Endpoint, hc *HealthChecker) *EndpointChecker {
	c := &EndpointChecker{
		tcpChecker:    newTcpChecker(endpoint, hc.timeout),
		endpoint:      endpoint,
		HealthChecker: hc,
		resp:          make(chan checkResponse),
		timeout:       make(chan bool),
		stop:          make(chan struct{}),
	}
	return c
}

func newTcpChecker(endpoint *model.Endpoint, timeout time.Duration) *TCPChecker {
	return &TCPChecker{
		addr:    endpoint.Address.GetAddress(),
		timeout: timeout,
	}
}

func (hc *HealthChecker) getCheckInterval() time.Duration {
	return hc.intervalBase
}

func (c *EndpointChecker) Start() {
	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("[health check] node checker panic %v\n%s", r, string(debug.Stack()))
		}
		c.checkTimer.Stop()
		c.checkTimeout.Stop()
	}()
	c.checkTimer = gxtime.AfterFunc(c.HealthChecker.initialDelay, c.OnCheck)
	for {
		select {
		case <-c.stop:
			return
		default:
			// prepare a check
			currentID := atomic.AddUint64(&c.checkID, 1)
			select {
			case <-c.stop:
				return
			case resp := <-c.resp:
				if resp.ID == currentID {
					if c.checkTimeout != nil {
						c.checkTimeout.Stop()
						c.checkTimeout = nil
					}
					if resp.Healthy {
						c.HandleSuccess()
					} else {
						c.HandleFailure(false)
					}
					c.checkTimer = gxtime.AfterFunc(c.HealthChecker.getCheckInterval(), c.OnCheck)
				}
			case <-c.timeout:
				if c.checkTimer != nil {
					c.checkTimer.Stop()
				}
				c.HandleFailure(true)
				c.checkTimer = gxtime.AfterFunc(c.HealthChecker.getCheckInterval(), c.OnCheck)
				logger.Infof("[health check] receive a timeout response at id: %d", currentID)

			}
		}
	}
}

func (c *EndpointChecker) Stop() {
	c.once.Do(func() {
		close(c.stop)
	})
}

func (c *EndpointChecker) HandleSuccess() {
	c.unHealthCount = 0
	c.healthCount++
	if c.healthCount > c.threshold {
		c.handleHealth()
	}
}

func (c *EndpointChecker) HandleFailure(timeout bool) {
	if timeout {
		c.HandleTimeout()
	} else {
		c.handleUnHealth()
	}
}

func (c *EndpointChecker) HandleTimeout() {
	c.healthCount = 0
	c.unHealthCount++
	if c.unHealthCount > c.threshold {
		c.handleUnHealth()
	}
}

func (c *EndpointChecker) handleHealth() {
	c.healthCount = 0
	c.unHealthCount = 0
	c.endpoint.UnHealthy = false
}

func (c *EndpointChecker) handleUnHealth() {
	c.healthCount = 0
	c.unHealthCount = 0
	c.endpoint.UnHealthy = true
}

func (c *EndpointChecker) OnCheck() {
	id := atomic.LoadUint64(&c.checkID)
	if c.checkTimeout != nil {
		c.checkTimeout.Stop()
	}
	c.checkTimeout = gxtime.AfterFunc(c.HealthChecker.timeout, c.OnTimeout)
	c.resp <- checkResponse{
		ID:      id,
		Healthy: c.tcpChecker.CheckHealth(),
	}
}

func (c *EndpointChecker) OnTimeout() {
	c.timeout <- true
}
