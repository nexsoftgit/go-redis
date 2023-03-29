package redis

import (
	"github.com/afex/hystrix-go/hystrix"
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

func generateMetric(command string, err error) {
	if err != nil {
		// circuitBreakerMetric.WithLabelValues(command, "go_redis_client", "ok", "").Inc()
		return
	}

	if err == hystrix.ErrCircuitOpen || err == hystrix.ErrMaxConcurrency || err == hystrix.ErrTimeout {
		circuitBreakerMetric.WithLabelValues(command, "go_redis_client", "fail", err.Error()).Inc()
	}

}
