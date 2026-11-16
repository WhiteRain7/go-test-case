package promanalytics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wallets_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"path", "method", "code"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "wallets_http_request_duration_seconds",
			Help:    "HTTP request latency.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func Init() {
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

type responseRecorder struct {
	http.ResponseWriter
	Code    int
	written bool
}

func (r *responseRecorder) WriteHeader(code int) {
	if !r.written {
		r.Code = code
		r.written = true
		r.ResponseWriter.WriteHeader(code)
	}
}

func MiddlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		record := &responseRecorder{res, 200, false}

		next.ServeHTTP(record, req)

		duration := time.Since(start).Seconds()
		path := req.URL.Path

		httpRequestsTotal.WithLabelValues(path, req.Method, strconv.Itoa(record.Code)).Inc()
		httpRequestDuration.WithLabelValues(path, req.Method).Observe(duration)
	})
}
