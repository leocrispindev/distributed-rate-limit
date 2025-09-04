package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricNameRequestsTotal = "requests_by_client"
	clientLabel             = "client"
)

type PrometheusRepository struct {
	counter *prometheus.CounterVec
}

func NewPrometheusRepository() *PrometheusRepository {
	requestsByClient := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: metricNameRequestsTotal,
			Help: "number of requests by client",
		},
		[]string{clientLabel}, //dynamic label (client name)
	)

	prometheus.MustRegister(requestsByClient)

	return &PrometheusRepository{
		counter: requestsByClient,
	}
}

func (r *PrometheusRepository) CountMetric(name string) error {
	r.counter.With(prometheus.Labels{
		clientLabel: name,
	}).Inc()

	return nil
}
