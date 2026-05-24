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
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

import (
	gxtime "github.com/dubbogo/gost/time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	DefaultTimeout  = time.Second
	DefaultInterval = 3 * time.Second

	DefaultHealthyThreshold   uint32 = 5
	DefaultUnhealthyThreshold uint32 = 5
	DefaultFirstInterval             = 5 * time.Second
)

type HealthChecker struct {
	// checkers is the live set of per-address health-check sessions.
	//
	// Concurrency contract: mutated only from ClusterStore mutation paths
	// (AddCluster, UpdateCluster, SetEndpoint, DeleteEndpoint, and
	// ensureRuntimeClusters), all of which hold ClusterManager.rw for write.
	// EndpointChecker goroutines never read or mutate this map; they only
	// touch their own captured fields. The map is therefore unlocked by
	// design. Do not call Start, Stop, StartOne, or StopOne from a goroutine
	// that does not hold the ClusterManager write lock.
	checkers      map[string]*EndpointChecker
	sessionConfig map[string]any
	// check config
	timeout            time.Duration
	intervalBase       time.Duration
	healthyThreshold   uint32
	initialDelay       time.Duration
	cluster            *model.ClusterConfig
	unhealthyThreshold uint32
	protocol           string
	onEndpointHealth   EndpointHealthListener
}

// EndpointChecker is a wrapper of types.HealthCheckSession for health check
type EndpointChecker struct {
	endpoint      *model.Endpoint
	endpointID    string
	endpointAddr  string
	HealthChecker *HealthChecker
	// checker, todo can extend to TCP, http, grpc, dubbo or other protocol checker
	checker       Checker
	resp          chan checkResponse
	timeout       chan bool
	checkID       uint64
	stop          chan struct{}
	checkTimer    *gxtime.Timer
	checkTimeout  *gxtime.Timer
	unHealthCount uint32
	healthCount   uint32

	once sync.Once
}

type Checker interface {
	CheckHealth() bool
	OnTimeout()
}

type EndpointHealthEvent struct {
	EndpointID      string
	EndpointAddress string
	Healthy         bool
}

type EndpointHealthListener func(EndpointHealthEvent)

type checkResponse struct {
	ID      uint64
	Healthy bool
}

func CreateHealthCheck(cluster *model.ClusterConfig, cfg model.HealthCheckConfig) *HealthChecker {
	return CreateHealthCheckWithCallback(cluster, cfg, nil)
}

func CreateHealthCheckWithCallback(
	cluster *model.ClusterConfig,
	cfg model.HealthCheckConfig,
	onEndpointHealth EndpointHealthListener,
) *HealthChecker {

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

	initialDelay := DefaultFirstInterval
	initialDelaySeconds, err := strconv.Atoi(cfg.InitialDelaySeconds)
	if err != nil {
		logger.Infof("[health check] initialDelay parse seconds error %s", err)
	} else {
		initialDelay = time.Duration(initialDelaySeconds) * time.Second
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
		protocol:           cfg.Protocol,
		sessionConfig:      cfg.SessionConfig,
		cluster:            cluster,
		timeout:            timeout,
		intervalBase:       interval,
		healthyThreshold:   healthyThreshold,
		unhealthyThreshold: unhealthyThreshold,
		initialDelay:       initialDelay,
		checkers:           make(map[string]*EndpointChecker),
		onEndpointHealth:   onEndpointHealth,
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
	for addr, h := range hc.checkers {
		h.Stop()
		delete(hc.checkers, addr)
		logger.Infof("[health check] stop a health check session for %s", addr)
	}
}

func (hc *HealthChecker) StopOne(endpoint *model.Endpoint) {
	hc.stopCheck(endpoint)
}

func (hc *HealthChecker) StartOne(endpoint *model.Endpoint) {
	hc.startCheck(endpoint)
}

func (hc *HealthChecker) startCheck(endpoint *model.Endpoint) {
	if endpoint == nil {
		return
	}
	addr := endpoint.Address.GetAddress()
	if _, ok := hc.checkers[addr]; !ok {
		c := newChecker(endpoint, hc)
		hc.checkers[addr] = c
		go c.Start()
		logger.Infof("[health check] create a health check session for %s", addr)
	}
}

func (hc *HealthChecker) stopCheck(endpoint *model.Endpoint) {
	if endpoint == nil {
		return
	}
	addr := endpoint.Address.GetAddress()
	if hc.hasOtherEndpointWithAddress(endpoint, addr) {
		return
	}
	if c, ok := hc.checkers[addr]; ok {
		c.Stop()
		delete(hc.checkers, addr)
		logger.Infof("[health check] stop a health check session for %s", addr)
	}
}

func (hc *HealthChecker) hasOtherEndpointWithAddress(endpoint *model.Endpoint, addr string) bool {
	if hc.cluster == nil {
		return false
	}
	for _, candidate := range hc.cluster.Endpoints {
		if candidate == nil || candidate == endpoint {
			continue
		}
		if endpoint.ID != "" && candidate.ID == endpoint.ID {
			continue
		}
		if candidate.Address.GetAddress() == addr {
			return true
		}
	}
	return false
}

func newChecker(endpoint *model.Endpoint, hc *HealthChecker) *EndpointChecker {
	var checker Checker
	protocol := strings.ToLower(hc.protocol)
	switch protocol {
	case "tcp":
		checker = &TCPChecker{
			address: endpoint.Address.GetAddress(),
			timeout: hc.timeout,
		}
	case "http":
		checker = &HTTPChecker{
			address: endpoint.Address.GetAddress(),
			timeout: hc.timeout,
		}
	case "https":
		checker = &HTTPSChecker{
			address: endpoint.Address.GetAddress(),
			timeout: hc.timeout,
		}
	default:
		logger.Warnf("[health check] %s health checker is not implemented, using tcp checker", hc.protocol)
		checker = &TCPChecker{
			address: endpoint.Address.GetAddress(),
			timeout: hc.timeout,
		}
	}

	c := &EndpointChecker{
		checker:  checker,
		endpoint: endpoint,
		// Capture stable event identity when the checker starts; endpoint
		// objects can be replaced while old checker events are still in flight.
		endpointID:    endpoint.ID,
		endpointAddr:  endpoint.Address.GetAddress(),
		HealthChecker: hc,
		resp:          make(chan checkResponse),
		timeout:       make(chan bool),
		stop:          make(chan struct{}),
	}
	return c
}

func (hc *HealthChecker) getCheckInterval() time.Duration {
	return hc.intervalBase
}

func (c *EndpointChecker) Start() {
	defer c.cleanupStart()
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

func (c *EndpointChecker) cleanupStart() {
	if r := recover(); r != nil {
		logger.Warnf("[health check] node checker panic %v\n%s", r, string(debug.Stack()))
	}
	c.stopCheckTimers()
}

func (c *EndpointChecker) stopCheckTimers() {
	// Early Stop may run before timers are assigned.
	if c.checkTimer != nil {
		c.checkTimer.Stop()
	}
	if c.checkTimeout != nil {
		c.checkTimeout.Stop()
	}
}

func (c *EndpointChecker) Stop() {
	c.once.Do(func() {
		close(c.stop)
	})
}

// HandleSuccess records a healthy probe result. The endpoint is only
// flipped to healthy after healthyThreshold consecutive successes — see
// the CHANGELOG for the v1.2 behavior change that made this configured
// threshold actually take effect (previously the counter compared
// against an uninitialized field that was always 0, so the first probe
// flipped state).
func (c *EndpointChecker) HandleSuccess() {
	c.unHealthCount = 0
	c.healthCount++
	if c.healthCount >= c.HealthChecker.healthyThreshold {
		c.handleHealth()
	}
}

// HandleFailure records an unhealthy probe result. Both negative probe
// responses and timeouts feed the same unhealthy counter, so the configured
// unhealthyThreshold governs the flip from healthy to unhealthy regardless of
// how the failure manifested. Prior to v1.2, timeout=false flipped state
// immediately and timeout=true compared against an uninitialized threshold
// field; see CHANGELOG.
//
// Deprecated: the timeout argument is ignored. It is retained only because
// HandleFailure was exported in v1.x and external Checker implementations may
// still pass it. New code should call this method with false; the next major
// release will collapse the signature to HandleFailure().
func (c *EndpointChecker) HandleFailure(timeout bool) {
	_ = timeout
	c.healthCount = 0
	c.unHealthCount++
	if c.unHealthCount >= c.HealthChecker.unhealthyThreshold {
		c.handleUnHealth()
	}
}

// HandleTimeout is preserved for backward compatibility with external Checker
// implementations that called it directly.
//
// Deprecated: routes to HandleFailure(true). Will be removed in the next major
// release.
func (c *EndpointChecker) HandleTimeout() {
	c.HandleFailure(true)
}

func (c *EndpointChecker) handleHealth() {
	c.healthCount = 0
	c.unHealthCount = 0
	c.emitHealth(true)
}

func (c *EndpointChecker) handleUnHealth() {
	c.healthCount = 0
	c.unHealthCount = 0
	c.emitHealth(false)
}

func (c *EndpointChecker) emitHealth(healthy bool) {
	if c.HealthChecker.onEndpointHealth == nil {
		// Direct CreateHealthCheck (no-callback) path. In-tree this branch is
		// unreachable: all in-tree HealthCheckers are constructed via
		// CreateHealthCheckWithCallback (see pkg/cluster/cluster.go). The
		// mutation below is preserved for external consumers that import
		// healthcheck directly. The cluster snapshot does not observe this
		// mutation; in-tree health flows go through the callback above.
		if !c.HealthChecker.setEndpointAddressHealth(c.endpointAddr, healthy) {
			c.endpoint.UnHealthy = !healthy
		}
		return
	}
	c.HealthChecker.onEndpointHealth(EndpointHealthEvent{
		EndpointID:      c.endpointID,
		EndpointAddress: c.endpointAddr,
		Healthy:         healthy,
	})
}

func (hc *HealthChecker) setEndpointAddressHealth(addr string, healthy bool) bool {
	if hc.cluster == nil {
		return false
	}
	updated := false
	for _, endpoint := range hc.cluster.Endpoints {
		if endpoint == nil || endpoint.Address.GetAddress() != addr {
			continue
		}
		endpoint.UnHealthy = !healthy
		updated = true
	}
	return updated
}

func (c *EndpointChecker) OnCheck() {
	id := atomic.LoadUint64(&c.checkID)
	if c.checkTimeout != nil {
		c.checkTimeout.Stop()
	}
	c.checkTimeout = gxtime.AfterFunc(c.HealthChecker.timeout, c.OnTimeout)
	c.resp <- checkResponse{
		ID:      id,
		Healthy: c.checker.CheckHealth(),
	}
}

func (c *EndpointChecker) OnTimeout() {
	c.timeout <- true
}
