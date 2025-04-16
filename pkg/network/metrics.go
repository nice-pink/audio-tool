package network

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// metrics

type Metrics struct {
	writeLoopMetric  prometheus.Counter
	writeBytesMetric prometheus.Counter
	readBytesMetric  prometheus.Counter
}

func NewMetrics(prefix string, labels map[string]string) *Metrics {
	m := &Metrics{}

	//
	m.writeLoopMetric = promauto.NewCounter(prometheus.CounterOpts{
		Subsystem:   prefix,
		Name:        "write_loop",
		Help:        "Write loop counter.",
		ConstLabels: prometheus.Labels(labels),
	})
	prometheus.MustRegister(m.writeLoopMetric)

	m.writeBytesMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "bytes_written",
		Help:        "Write loop counter.",
		ConstLabels: prometheus.Labels(labels),
	})
	prometheus.MustRegister(m.writeBytesMetric)

	m.readBytesMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "bytes_read",
		Help:        "Read loop counter.",
		ConstLabels: prometheus.Labels(labels),
	})
	prometheus.MustRegister(m.readBytesMetric)

	return m
}
