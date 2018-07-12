package redis

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	circuitBreakerOpen = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_open",
			Help: "A total number of circuit breaker state open. This happens due to the circuit being measured as unhealthy.",
		}, []string{"action"})

	circuitBreakerMaxConcurrency = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_max_concurrency",
			Help: "A total number of client executed at the same time and exceeded max concurrency.",
		}, []string{"action"})

	circuitBreakerTimeout = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "circuit_breaker_timeout",
			Help: "A total number of request exceeded timeout duration",
		}, []string{"action"})
)

func init() {
	// Metrics have to be registerd to be exposed
	prometheus.MustRegister(
		circuitBreakerMaxConcurrency,
		circuitBreakerOpen,
		circuitBreakerTimeout)
}
