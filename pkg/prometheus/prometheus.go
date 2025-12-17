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

package prometheus

import (
	"bytes"
	"errors"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/common/expfmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/context"
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var defaultSubsystem = "pixiu"

type ContextHandlerFunc func(c *contextHttp.HttpContext) error

const (
	_          = iota // ignore first value by assigning to blank identifier
	KB float64 = 1 << (10 * iota)
	MB
	GB
	TB
)

type FavContextKeyType string

type Metric struct {
	MetricCollector prometheus.Collector
	ID              string
	Name            string
	Description     string
	Type            string
	Args            []string
	Buckets         []float64
}

// reqDurBuckets is the buckets for request duration. Here, we use the prometheus defaults
// which are for ~10s request length max: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}
var reqDurBuckets = prometheus.DefBuckets

// reqSzBuckets is the buckets for request size. Here we define a spectrom from 1KB thru 1NB up to 10MB.
var reqSzBuckets = []float64{1.0 * KB, 2.0 * KB, 5.0 * KB, 10.0 * KB, 100 * KB, 500 * KB, 1.0 * MB, 2.5 * MB, 5.0 * MB, 10.0 * MB}

// resSzBuckets is the buckets for response size. Here we define a spectrom from 1KB thru 1NB up to 10MB.
var resSzBuckets = []float64{1.0 * KB, 2.0 * KB, 5.0 * KB, 10.0 * KB, 100 * KB, 500 * KB, 1.0 * MB, 2.5 * MB, 5.0 * MB, 10.0 * MB}

//  Standard default metrics
//	counter, counter_vec, gauge, gauge_vec,
//	histogram, histogram_vec, summary, summary_vec

// Deprecated: Use reqCntNew instead. This will be removed in future versions.
var reqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total",
	Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "host", "url"},
}

var reqCntNew = &Metric{
	ID:          "reqCntNew",
	Name:        "request_count",
	Description: "request total count in pixiu",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "host", "url"},
}

var reqElapsed = &Metric{
	ID:          "reqElapsed",
	Name:        "request_elapsed",
	Description: "request total elapsed in pixiu (milliseconds)",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "host", "url"},
}

var reqErrorCnt = &Metric{
	ID:          "reqErrorCnt",
	Name:        "request_error_count",
	Description: "request error total count in pixiu",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "host", "url"},
}

var reqDur = &Metric{
	ID:          "reqDur",
	Name:        "process_time_millisec",
	Description: "request process time response in pixiu (milliseconds)",
	Args:        []string{"code", "method", "url"},
	Type:        "histogram_vec",
	Buckets:     reqDurBuckets,
}

var resSz = &Metric{
	ID:          "resSz",
	Name:        "response_content_length",
	Description: "request total content length response in pixiu (bytes)",
	Args:        []string{"code", "method", "url"},
	Type:        "counter_vec",
}

var reqSz = &Metric{
	ID:          "reqSz",
	Name:        "request_content_length",
	Description: "request total content length in pixiu (bytes)",
	Args:        []string{"code", "method", "url"},
	Type:        "counter_vec",
}

var standardMetrics = []*Metric{
	reqCnt,    // Deprecated: for backward compatibility
	reqCntNew, // New unified metric name
	reqElapsed,
	reqErrorCnt,
	reqDur,
	resSz,
	reqSz,
}

// NewMetric associates prometheus.Collector based on Metric.Type
func NewMetric(m *Metric, subsystem string) prometheus.Collector {
	var metric prometheus.Collector
	switch m.Type {
	case "counter_vec":
		metric = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "counter":
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "gauge_vec":
		metric = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "gauge":
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	case "histogram_vec":
		metric = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
				Buckets:   m.Buckets,
			},
			m.Args,
		)
	case "histogram":
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
				Buckets:   m.Buckets,
			},
		)
	case "summary_vec":
		metric = prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
			m.Args,
		)
	case "summary":
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Subsystem: subsystem,
				Name:      m.Name,
				Help:      m.Description,
			},
		)
	}
	return metric
}

type RequestCounterLabelMappingFunc func(c *contextHttp.HttpContext) string

type Prometheus struct {
	reqCnt       *prometheus.CounterVec // Deprecated
	reqCntNew    *prometheus.CounterVec // New unified name
	reqElapsed   *prometheus.CounterVec
	reqErrorCnt  *prometheus.CounterVec
	reqDur       *prometheus.HistogramVec
	reqSz, resSz *prometheus.CounterVec
	Ppg          PushGateway

	MetricsList []*Metric
	MetricsPath string
	Subsystem   string

	RequestCounterURLLabelMappingFunc  RequestCounterLabelMappingFunc
	RequestCounterHostLabelMappingFunc RequestCounterLabelMappingFunc

	URLLabelFromContext string
	Datacontext         context.Context

	// Dynamic metrics storage for custom metrics
	dynamicCounters   sync.Map // map[string]*prometheus.CounterVec
	dynamicGauges     sync.Map // map[string]*prometheus.GaugeVec
	dynamicHistograms sync.Map // map[string]*prometheus.HistogramVec
}

// PushGateway contains the configuration for pushing to a Prometheus pushgateway (optional)
type PushGateway struct {
	CounterPush           bool
	PushInterval          time.Duration
	PushIntervalThreshold int
	PushGatewayURL        string
	Job                   string
	counter               int
	mutex                 sync.RWMutex
}

// NewPrometheus generates a new set of metrics with a certain subsystem name
func NewPrometheus() *Prometheus {
	var metricsList []*Metric
	metricsList = append(metricsList, standardMetrics...)
	p := &Prometheus{
		MetricsList: metricsList,
		Subsystem:   defaultSubsystem,
		RequestCounterURLLabelMappingFunc: func(c *contextHttp.HttpContext) string {
			return c.GetUrl()
		},
		RequestCounterHostLabelMappingFunc: func(c *contextHttp.HttpContext) string {
			return c.Request.Host
		},
	}
	p.registerMetrics()
	return p
}

func (p *Prometheus) registerMetrics() {
	for _, metricDef := range p.MetricsList {
		metric := NewMetric(metricDef, p.Subsystem)
		if err := prometheus.Register(metric); err != nil {
			logger.Errorf("%s could not be registered in Prometheus: %v", metricDef.Name, err)
		}
		switch metricDef {

		case reqCnt:
			p.reqCnt = metric.(*prometheus.CounterVec)
		case reqCntNew:
			p.reqCntNew = metric.(*prometheus.CounterVec)
		case reqElapsed:
			p.reqElapsed = metric.(*prometheus.CounterVec)
		case reqErrorCnt:
			p.reqErrorCnt = metric.(*prometheus.CounterVec)
		case reqDur:
			p.reqDur = metric.(*prometheus.HistogramVec)
		case resSz:
			p.resSz = metric.(*prometheus.CounterVec)
		case reqSz:
			p.reqSz = metric.(*prometheus.CounterVec)
		}
		metricDef.MetricCollector = metric
	}
}

func (p *Prometheus) SetPushGatewayUrl(pushGatewayURL, metricspath string) {
	p.Ppg.mutex.Lock()
	defer p.Ppg.mutex.Unlock()
	p.Ppg.PushGatewayURL = pushGatewayURL
	p.MetricsPath = metricspath
}

func (p *Prometheus) SetPushIntervalThreshold(isTurn bool, pushIntervalThreshold int) {
	p.Ppg.CounterPush = isTurn
	p.Ppg.PushIntervalThreshold = pushIntervalThreshold
}

func (p *Prometheus) SetPushGatewayJob(j string) {
	p.Ppg.mutex.Lock()
	defer p.Ppg.mutex.Unlock()
	p.Ppg.Job = j
}

func (p *Prometheus) startPushCounter() {
	if p.Ppg.counter >= p.Ppg.PushIntervalThreshold {
		go p.sendMetricsToPushGateway(p.getMetrics())
		p.Ppg.counter = 0
	}
}

func (p *Prometheus) SetPushGateway() {
	if p.Ppg.CounterPush {
		p.startPushCounter()
	}
}

func (p *Prometheus) getMetrics() []byte {
	out := &bytes.Buffer{}
	metricFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		logger.Errorf("prometheus.DefaultGatherer.Gather error: %v", err)
		return []byte{}
	}
	for i := range metricFamilies {
		_, err := expfmt.MetricFamilyToText(out, metricFamilies[i])
		if err != nil {
			logger.Errorf("failed to converts a MetricFamily proto message into text format %v", err)
		}
	}
	return out.Bytes()
}

func (p *Prometheus) sendMetricsToPushGateway(metrics []byte) {
	req, err := http.NewRequest(http.MethodPost, p.getPushGatewayURL(), bytes.NewBuffer(metrics))
	if err != nil {
		logger.Errorf("failed to create push gateway request: %v", err)
		return
	}
	if _, err = (&http.Client{}).Do(req); err != nil {
		logger.Errorf("Error sending to push gateway: %v", err)
	}
}

func (p *Prometheus) getPushGatewayURL() string {
	p.Ppg.mutex.Lock()
	defer p.Ppg.mutex.Unlock()
	h, _ := os.Hostname()
	if p.Ppg.Job == "" {
		p.Ppg.Job = "pixiu"
	}
	return p.Ppg.PushGatewayURL + p.MetricsPath + "/job/" + p.Ppg.Job + "/instance/" + h
}

// HandlerFunc defines handler function for middleware
func (p *Prometheus) HandlerFunc() ContextHandlerFunc {
	return func(c *contextHttp.HttpContext) error {
		start := time.Now()
		reqSz, err1 := computeApproximateRequestSize(c.Request)
		//fmt.Println("reqSz", reqSz)
		elapsed := float64(time.Since(start).Milliseconds())
		//fmt.Println("elapsed ", elapsed)
		url := p.RequestCounterURLLabelMappingFunc(c)
		//fmt.Println("url ", url)
		statusStr := strconv.Itoa(c.GetStatusCode())
		//fmt.Println("statusStr", statusStr)
		method := c.GetMethod()
		//fmt.Println("method ", method)
		host := p.RequestCounterHostLabelMappingFunc(c)

		// Record metrics aligned with Pull mode
		// Update both old (deprecated) and new metric names for backward compatibility
		p.reqCnt.WithLabelValues(statusStr, method, host, url).Inc()    // Deprecated: will be removed
		p.reqCntNew.WithLabelValues(statusStr, method, host, url).Inc() // New unified name
		p.reqElapsed.WithLabelValues(statusStr, method, host, url).Add(elapsed)
		p.reqDur.WithLabelValues(statusStr, method, url).Observe(elapsed)

		if err1 == nil {
			p.reqSz.WithLabelValues(statusStr, method, url).Add(float64(reqSz))
		}
		resSz, err2 := computeApproximateResponseSize(c.TargetResp)
		if err2 == nil {
			p.resSz.WithLabelValues(statusStr, method, url).Add(float64(resSz))
		}

		// Record errors
		if c.LocalReply() {
			p.reqErrorCnt.WithLabelValues(statusStr, method, host, url).Inc()
		}

		p.Ppg.mutex.Lock()
		p.Ppg.counter = p.Ppg.counter + 1
		defer p.Ppg.mutex.Unlock()
		p.SetPushGateway()
		return nil
	}
}

func computeApproximateRequestSize(r *http.Request) (int, error) {
	if r == nil {
		return 0, errors.New("http.Request is null pointer ")
	}
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}
	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)
	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s, nil
}

func computeApproximateResponseSize(res any) (int, error) {
	if res == nil {
		return 0, errors.New("client response is nil")
	}
	if unaryResponse, ok := res.(*client.UnaryResponse); ok {
		return len(unaryResponse.Data), nil
	}
	return 0, errors.New("response is not of type client.UnaryResponse")
}

// RecordDynamicMetric records a dynamic metric based on type (counter, gauge, histogram)
func (p *Prometheus) RecordDynamicMetric(name string, metricType string, value float64, labels map[string]string) error {
	// Extract label keys and values
	labelKeys := make([]string, 0, len(labels))
	labelValues := make([]string, 0, len(labels))
	for k, v := range labels {
		labelKeys = append(labelKeys, k)
		labelValues = append(labelValues, v)
	}

	switch metricType {
	case "counter":
		return p.recordDynamicCounter(name, value, labelKeys, labelValues)
	case "gauge":
		return p.recordDynamicGauge(name, value, labelKeys, labelValues)
	case "histogram":
		return p.recordDynamicHistogram(name, value, labelKeys, labelValues)
	default:
		return errors.New("unsupported metric type: " + metricType)
	}
}

// recordDynamicCounter records a counter metric
func (p *Prometheus) recordDynamicCounter(name string, value float64, labelKeys, labelValues []string) error {
	// Create a unique key for this metric with its label keys
	metricKey := name + "_" + joinLabels(labelKeys)

	// Try to load existing counter
	if metric, ok := p.dynamicCounters.Load(metricKey); ok {
		counter := metric.(*prometheus.CounterVec)
		counter.WithLabelValues(labelValues...).Add(value)
		return nil
	}

	// Create new counter
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: p.Subsystem,
			Name:      name,
			Help:      "Dynamic counter: " + name,
		},
		labelKeys,
	)

	// Register the metric
	if err := prometheus.Register(counter); err != nil {
		// Metric might already be registered, try to use it
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			counter = are.ExistingCollector.(*prometheus.CounterVec)
		} else {
			return err
		}
	}

	// Store for future use
	p.dynamicCounters.Store(metricKey, counter)
	counter.WithLabelValues(labelValues...).Add(value)
	return nil
}

// recordDynamicGauge records a gauge metric
func (p *Prometheus) recordDynamicGauge(name string, value float64, labelKeys, labelValues []string) error {
	// Create a unique key for this metric with its label keys
	metricKey := name + "_" + joinLabels(labelKeys)

	// Try to load existing gauge
	if metric, ok := p.dynamicGauges.Load(metricKey); ok {
		gauge := metric.(*prometheus.GaugeVec)
		gauge.WithLabelValues(labelValues...).Set(value)
		return nil
	}

	// Create new gauge
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: p.Subsystem,
			Name:      name,
			Help:      "Dynamic gauge: " + name,
		},
		labelKeys,
	)

	// Register the metric
	if err := prometheus.Register(gauge); err != nil {
		// Metric might already be registered, try to use it
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			gauge = are.ExistingCollector.(*prometheus.GaugeVec)
		} else {
			return err
		}
	}

	// Store for future use
	p.dynamicGauges.Store(metricKey, gauge)
	gauge.WithLabelValues(labelValues...).Set(value)
	return nil
}

// recordDynamicHistogram records a histogram metric
func (p *Prometheus) recordDynamicHistogram(name string, value float64, labelKeys, labelValues []string) error {
	// Create a unique key for this metric with its label keys
	metricKey := name + "_" + joinLabels(labelKeys)

	// Try to load existing histogram
	if metric, ok := p.dynamicHistograms.Load(metricKey); ok {
		histogram := metric.(*prometheus.HistogramVec)
		histogram.WithLabelValues(labelValues...).Observe(value)
		return nil
	}

	// Create new histogram with default buckets
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Subsystem: p.Subsystem,
			Name:      name,
			Help:      "Dynamic histogram: " + name,
			Buckets:   prometheus.DefBuckets,
		},
		labelKeys,
	)

	// Register the metric
	if err := prometheus.Register(histogram); err != nil {
		// Metric might already be registered, try to use it
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			histogram = are.ExistingCollector.(*prometheus.HistogramVec)
		} else {
			return err
		}
	}

	// Store for future use
	p.dynamicHistograms.Store(metricKey, histogram)
	histogram.WithLabelValues(labelValues...).Observe(value)
	return nil
}

// joinLabels creates a consistent key from label keys
func joinLabels(labels []string) string {
	if len(labels) == 0 {
		return "no_labels"
	}
	result := ""
	for i, label := range labels {
		if i > 0 {
			result += "_"
		}
		result += label
	}
	return result
}
