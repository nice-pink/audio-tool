package encodings

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	validationErrorMetric prometheus.Counter
	audioParseErrorMetric prometheus.Counter
}

func NewMetrics(prefix string, labels map[string]string) *Metrics {
	m := &Metrics{}

	m.validationErrorMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "validation_error",
		Help:        "Validation error counter.",
		ConstLabels: prometheus.Labels(labels),
	})
	prometheus.MustRegister(m.validationErrorMetric)

	m.audioParseErrorMetric = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "parse_audio_error",
		Help:        "Validation error counter.",
		ConstLabels: prometheus.Labels(labels),
	})
	prometheus.MustRegister(m.audioParseErrorMetric)

	return m
}
