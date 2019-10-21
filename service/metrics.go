package service

import (
	"github.com/integration-system/isp-lib/metric"
	"github.com/rcrowley/go-metrics"
	"strconv"
	"sync"
	"time"
)

const (
	defaultSampleSize = 2048
)

var Metrics = metricService{mh: nil}

type (
	metricService struct {
		mh *metricHolder
	}

	metricHolder struct {
		methodHistograms   map[string]metrics.Histogram
		methodLock         sync.RWMutex
		statusCounters     map[int]metrics.Counter
		statusLock         sync.RWMutex
		routerResponseTime metrics.Histogram
		responseTime       metrics.Histogram
	}
)

func (m metricService) Init() {
	if m.mh == nil {
		m.mh = &metricHolder{
			methodHistograms: make(map[string]metrics.Histogram),
			statusCounters:   make(map[int]metrics.Counter),
			responseTime: metrics.GetOrRegisterHistogram(
				"http.response.time", metric.GetRegistry(), metrics.NewUniformSample(defaultSampleSize),
			),
			routerResponseTime: metrics.GetOrRegisterHistogram(
				"grpc.router.response.time", metric.GetRegistry(), metrics.NewUniformSample(defaultSampleSize),
			),
		}
	}
}

func (m metricService) UpdateMethodResponseTime(uri string, time time.Duration) {
	m.getOrRegisterHistogram(uri).Update(int64(time))
}

func (m metricService) UpdateResponseTime(time time.Duration) {
	m.mh.responseTime.Update(int64(time))
}

func (m metricService) UpdateRouterResponseTime(time time.Duration) {
	m.mh.routerResponseTime.Update(int64(time))
}

func (m metricService) UpdateStatusCounter(status int) {
	m.getOrRegisterCounter(status).Inc(1)
}

func (m metricService) getOrRegisterHistogram(uri string) metrics.Histogram {
	m.mh.methodLock.RLock()
	histogram, ok := m.mh.methodHistograms[uri]
	m.mh.methodLock.RUnlock()
	if ok {
		return histogram
	}

	m.mh.methodLock.Lock()
	defer m.mh.methodLock.Unlock()
	if d, ok := m.mh.methodHistograms[uri]; ok {
		return d
	}
	histogram = metrics.GetOrRegisterHistogram(
		"http.response.time_"+uri,
		metric.GetRegistry(),
		metrics.NewUniformSample(defaultSampleSize),
	)
	m.mh.methodHistograms[uri] = histogram
	return histogram
}

func (m metricService) getOrRegisterCounter(status int) metrics.Counter {
	m.mh.statusLock.RLock()
	d, ok := m.mh.statusCounters[status]
	m.mh.statusLock.RUnlock()
	if ok {
		return d
	}

	m.mh.statusLock.Lock()
	defer m.mh.statusLock.Unlock()
	if d, ok := m.mh.statusCounters[status]; ok {
		return d
	}
	d = metrics.GetOrRegisterCounter("http.response.count."+strconv.Itoa(status), metric.GetRegistry())
	m.mh.statusCounters[status] = d
	return d
}
