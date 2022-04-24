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
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	gxtime "github.com/dubbogo/gost/time"
	"runtime/debug"
	"sync/atomic"
	"time"
)

const (
	DefaultTimeout        = time.Second
	DefaultInterval       = 15 * time.Second
	DefaultIntervalJitter = 5 * time.Millisecond

	DefaultHealthyThreshold   uint32 = 1
	DefaultUnhealthyThreshold uint32 = 1
	DefaultFirstInterval             = time.Second
)

type HealthChecker struct {
	checkers      map[string]*EndpointChecker
	sessionConfig map[string]interface{}
	// check config
	timeout             time.Duration
	intervalBase        time.Duration
	intervalJitter      time.Duration
	healthyThreshold    uint32
	initialDelay        time.Duration
	cluster             *model.ClusterConfig
	unhealthyThreshold  uint32
	localProcessHealthy int64
}

// EndpointChecker is a wrapper of types.HealthCheckSession for health check
type EndpointChecker struct {
	endpoint      *model.Endpoint
	HealthChecker *HealthChecker
	//
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
}

type checkResponse struct {
	ID      uint64
	Healthy bool
}

func CreateHealthCheck(cluster *model.ClusterConfig, cfg model.HealthCheckConfig) *HealthChecker {
	timeout := cfg.TimeoutConfig
	if cfg.TimeoutConfig == 0 {
		timeout = DefaultTimeout
	}
	interval := cfg.IntervalConfig
	if cfg.IntervalConfig == 0 {
		interval = DefaultInterval
	}
	unhealthyThreshold := cfg.UnhealthyThreshold
	if unhealthyThreshold == 0 {
		unhealthyThreshold = DefaultUnhealthyThreshold
	}
	healthyThreshold := cfg.HealthyThreshold
	if healthyThreshold == 0 {
		healthyThreshold = DefaultHealthyThreshold
	}
	intervalJitter := cfg.IntervalJitterConfig
	if intervalJitter == 0 {
		intervalJitter = DefaultIntervalJitter
	}
	initialDelay := DefaultFirstInterval
	if cfg.InitialDelaySeconds.Microseconds() > 0 {
		initialDelay = cfg.InitialDelaySeconds
	}

	hc := &HealthChecker{
		// cfg
		sessionConfig:      cfg.SessionConfig,
		cluster:            cluster,
		timeout:            timeout,
		intervalBase:       interval,
		intervalJitter:     intervalJitter,
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

func (hc *HealthChecker) startCheck(endpoint *model.Endpoint) {
	addr := endpoint.Address.GetAddress()
	if _, ok := hc.checkers[addr]; !ok {
		c := newChecker(endpoint, hc)
		hc.checkers[addr] = c
		c.Start()
		atomic.AddInt64(&hc.localProcessHealthy, 1) // default host is healthy
		logger.Infof("[health check] create a health check session for %s", addr)
	}
}

func newChecker(endpoint *model.Endpoint, hc *HealthChecker) *EndpointChecker {
	c := &EndpointChecker{
		tcpChecker:    newTcpChecker(endpoint),
		endpoint:      endpoint,
		HealthChecker: hc,
		resp:          make(chan checkResponse),
		timeout:       make(chan bool),
		stop:          make(chan struct{}),
	}
	return c
}

func newTcpChecker(endpoint *model.Endpoint) *TCPChecker {
	return &TCPChecker{
		addr: endpoint.Address.GetAddress(),
	}
}

func (hc *HealthChecker) getCheckInterval() time.Duration {
	interval := hc.intervalBase
	return interval
}

func (c *EndpointChecker) Start() {
	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("[health check] node checker panic %v\n%s", r, string(debug.Stack()))
		}
		// stop all the timer when start is finished
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
					c.checkTimeout.Stop()
					if resp.Healthy {
						c.HandleSuccess()
					} else {
						c.HandleFailure(false)
					}
					// next health checker
					c.checkTimer = gxtime.AfterFunc(c.HealthChecker.getCheckInterval(), c.OnCheck)
				}
			case <-c.timeout:
				c.checkTimer.Stop()
				c.HandleFailure(true)
				// next health checker
				c.checkTimer = gxtime.AfterFunc(c.HealthChecker.getCheckInterval(), c.OnCheck)
				logger.Infof("[health check] receive a timeout response at id: %d", currentID)

			}
		}
	}
}

func (c *EndpointChecker) Stop() {
	close(c.stop)
}

func (c *EndpointChecker) HandleSuccess() {
	c.unHealthCount = 0
	c.healthCount++
	if c.healthCount > c.threshold {
		c.handleHealth()
	}
}

func (c *EndpointChecker) HandleFailure(isTimeout bool) {
	if isTimeout {
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
	c.endpoint.Healthy = true
}

func (c *EndpointChecker) handleUnHealth() {
	c.healthCount = 0
	c.unHealthCount = 0
	c.endpoint.Healthy = false
}

func (c *EndpointChecker) OnCheck() {
	// record current id
	id := atomic.LoadUint64(&c.checkID)
	// start a timeout before check health
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
