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
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"
)

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	contextHttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
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

var reqCnt = &Metric{
	ID:          "reqCnt",
	Name:        "requests_total",
	Description: "How many HTTP requests processed, partitioned by status code and HTTP method.",
	Type:        "counter_vec",
	Args:        []string{"code", "method", "host", "url"},
}

var reqDur = &Metric{
	ID:          "reqDur",
	Name:        "request_duration_seconds",
	Description: "The HTTP request latencies in seconds.",
	Args:        []string{"code", "method", "url"},
	Type:        "histogram_vec",
	Buckets:     reqDurBuckets,
}

var resSz = &Metric{
	ID:          "resSz",
	Name:        "response_size_bytes",
	Description: "The HTTP response sizes in bytes.",
	Args:        []string{"code", "method", "url"},
	Type:        "histogram_vec",
	Buckets:     resSzBuckets,
}

var reqSz = &Metric{
	ID:          "reqSz",
	Name:        "request_size_bytes",
	Description: "The HTTP request sizes in bytes.",
	Args:        []string{"code", "method", "url"},
	Type:        "histogram_vec",
	Buckets:     reqSzBuckets,
}

var standardMetrics = []*Metric{
	reqCnt,
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
	reqCnt               *prometheus.CounterVec
	reqDur, reqSz, resSz *prometheus.HistogramVec
	Ppg                  PushGateway

	MetricsList []*Metric
	MetricsPath string
	Subsystem   string

	RequestCounterURLLabelMappingFunc  RequestCounterLabelMappingFunc
	RequestCounterHostLabelMappingFunc RequestCounterLabelMappingFunc

	URLLabelFromContext string
	Datacontext         context.Context
}

// PushGateway contains the configuration for pushing to a Prometheus pushgateway (optional)
type PushGateway struct {
	PushInterval   time.Duration
	PushGatewayURL string
	Job            string
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
		case reqDur:
			p.reqDur = metric.(*prometheus.HistogramVec)
		case resSz:
			p.resSz = metric.(*prometheus.HistogramVec)
		case reqSz:
			p.reqSz = metric.(*prometheus.HistogramVec)
		}
		metricDef.MetricCollector = metric
	}
}

func (p *Prometheus) SetMetricPath(path string) {
	p.MetricsPath = path
}

func (p *Prometheus) SetPushGatewayUrl(pushGatewayURL, metricsURL string, pushInterval time.Duration) {
	p.Ppg.PushGatewayURL = pushGatewayURL
	p.MetricsPath = metricsURL
	p.Ppg.PushInterval = pushInterval
	p.MetricsPath = metricsURL

}

func (p *Prometheus) SetPushGateway() {
	p.startPushTicker()
}

func (p *Prometheus) SetPushGatewayJob(j string) {
	p.Ppg.Job = j
}

func (p *Prometheus) startPushTicker() {
	p.sendMetricsToPushGateway(p.getMetrics())
}

func (p *Prometheus) getMetrics() []byte {

	out := &bytes.Buffer{}
	metricFamilies, _ := prometheus.DefaultGatherer.Gather()
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
	h, _ := os.Hostname()
	if p.Ppg.Job == "" {
		p.Ppg.Job = "pixiu"
	}
	return p.Ppg.PushGatewayURL + "/metrics/job/" + p.Ppg.Job + "/instance/" + h
}

// HandlerFunc defines handler function for middleware
func (p *Prometheus) HandlerFunc() ContextHandlerFunc {

	return func(c *contextHttp.HttpContext) error {

		start := time.Now()
		reqSz, err1 := computeApproximateRequestSize(c.Request)
		//fmt.Println("reqSz", reqSz)
		elapsed := float64(time.Since(start)) / float64(time.Second)
		//fmt.Println("elapsed ", elapsed)
		url := p.RequestCounterURLLabelMappingFunc(c)
		//fmt.Println("url ", url)
		statusStr := strconv.Itoa(c.GetStatusCode())
		//fmt.Println("statusStr", statusStr)
		method := c.GetMethod()
		//fmt.Println("method ", method)
		p.reqDur.WithLabelValues(statusStr, method, url).Observe(elapsed)
		p.reqCnt.WithLabelValues(statusStr, method, p.RequestCounterHostLabelMappingFunc(c), url).Inc()
		if err1 == nil {
			p.reqSz.WithLabelValues(statusStr, method, url).Observe(float64(reqSz))
		}
		resSz, err2 := computeApproximateResponseSize(c.TargetResp)
		if err2 == nil {
			p.resSz.WithLabelValues(statusStr, method, url).Observe(float64(resSz))
		}
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

func computeApproximateResponseSize(res *client.Response) (int, error) {
	if res == nil {
		return 0, errors.New("client.Response is null pointer ")
	}
	s := 0
	s += len(res.Data)
	return s, nil
}
