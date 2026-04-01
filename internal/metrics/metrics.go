package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"method", "endpoint"},
	)

	TransfersTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "transfers_total",
			Help: "Total number of successful transfers",
		},
	)

	TransfersFailed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "transfers_failed_total",
			Help: "Total number of failed transfers by reason",
		},
		[]string{"reason"},
	)

	TransferAmountTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "transfer_amount_paise_total",
			Help: "Total amount transferred in paise",
		},
	)
)
