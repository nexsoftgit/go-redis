package redis

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	circuitBreakerMetric = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker",
			Help: "A total number of circuit breaker state open. This happens due to the circuit being measured as unhealthy.",
		}, []string{"command", "service", "status", "state"})
)

func init() {
	// Metrics have to be registerd to be exposed
	prometheus.MustRegister(circuitBreakerMetric)
}
