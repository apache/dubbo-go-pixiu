package healthcheck

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	gxtime "github.com/dubbogo/gost/time"
	"net"
	"runtime/debug"
	"sync/atomic"
	"time"
)

// healthChecker is a basic implementation of a health checker.
// we use different implementations of types.Session to implement different health checker
type healthChecker struct {
	//
	checkers      map[string]*sessionChecker
	sessionConfig map[string]interface{}
	// check config
	timeout             time.Duration
	intervalBase        time.Duration
	intervalJitter      time.Duration
	healthyThreshold    uint32
	initialDelay        time.Duration
	cluster             *model.Cluster
	unhealthyThreshold  uint32
	localProcessHealthy int64
}

// EndpointChecker is a wrapper of types.HealthCheckSession for health check
type EndpointChecker struct {
	endpoint      *model.Endpoint
	HealthChecker *healthChecker
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
}

type checkResponse struct {
	ID      uint64
	Healthy bool
}

func (hc *healthChecker) start() {
	// each endpoint
	for _, h := range hc.cluster.Endpoints {
		hc.startCheck(h)
	}
}

func (hc *healthChecker) startCheck(endpoint *model.Endpoint) {
	addr := endpoint.Address.GetAddress()
	if _, ok := hc.checkers[addr]; !ok {
		c := newChecker(endpoint, hc)
		hc.checkers[addr] = c
		c.Start()
		atomic.AddInt64(&hc.localProcessHealthy, 1) // default host is healthy
		logger.Infof("[health check] create a health check session for %s", addr)
	}
}

func newChecker(endpoint *model.Endpoint, hc *healthChecker) *EndpointChecker {
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
				// if the ID is not equal, means we receive a timeout for this ID, ignore the response
				if resp.ID == currentID {
					c.checkTimeout.Stop()
					if resp.Healthy {
						c.HandleSuccess()
					} else {
						c.HandleFailure(types.FailureActive)
					}
					// next health checker
					c.checkTimer = gxtime.AfterFunc(c.HealthChecker.getCheckInterval(), c.OnCheck)
				}
			case <-c.timeout:
				c.checkTimer.Stop()
				c.HandleFailure(types.FailureNetwork)
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
	changed := false
	if c.endpoint.ContainHealthFlag(api.FAILED_ACTIVE_HC) {
		c.healthCount++
		// check the threshold
		if c.healthCount == c.HealthChecker.healthyThreshold {
			changed = true
			c.Host.ClearHealthFlag(api.FAILED_ACTIVE_HC)
		}
	}
	c.HealthChecker.incHealthy(c.Host, changed)
}

func (c *EndpointChecker) HandleFailure(reason types.FailureType) {
	c.healthCount = 0
	changed := false
	if !c.Host.ContainHealthFlag(api.FAILED_ACTIVE_HC) {
		c.unHealthCount++
		// check the threshold
		if c.unHealthCount == c.HealthChecker.unhealthyThreshold {
			changed = true
			c.Host.SetHealthFlag(api.FAILED_ACTIVE_HC)
		}
	}
	c.HealthChecker.decHealthy(c.Host, reason, changed)
}

func (c *EndpointChecker) OnCheck() {
	// record current id
	id := atomic.LoadUint64(&c.checkID)
	c.HealthChecker.stats.attempt.Inc(1)
	// start a timeout before check health
	c.checkTimeout.Stop()
	c.checkTimeout = utils.NewTimer(c.HealthChecker.timeout, c.OnTimeout)
	c.resp <- checkResponse{
		ID:      id,
		Healthy: c.Session.CheckHealth(),
	}
}

func (c *EndpointChecker) OnTimeout() {
	c.timeout <- true
}

type TCPChecker struct {
	addr string
}

func (s *TCPChecker) CheckHealth() bool {
	// default dial timeout, maybe already timeout by checker
	conn, err := net.DialTimeout("tcp", s.addr, 30*time.Second)
	if err != nil {
		logger.Error("[health check] tcp checker for host %s error: %v", s.addr, err)
		return false
	}
	conn.Close()
	return true
}

func (s *TCPChecker) OnTimeout() {}
