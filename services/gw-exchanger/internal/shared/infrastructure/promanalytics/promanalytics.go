package promanalytics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	Port int `mapstructure:"port"`
}

func ProvideMetrics(cfg Config) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.Port), nil)
}
