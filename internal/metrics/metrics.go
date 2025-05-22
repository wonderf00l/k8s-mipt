package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Время обработки запроса",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	LogRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "log_requests_total",
			Help: "Число /log запросов",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(RequestDuration, LogRequestsTotal)
}
