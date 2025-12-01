package main

import (
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Request duration",
		},
		[]string{"method", "path"},
	)

	buildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "build_info",
			Help: "Build information for GoBlog API",
		},
		[]string{"version", "build_date", "go_version", "env"},
	)
)

func init() {
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(buildInfo)
}

func SetBuildInfo(version, buildDate, env string) {
	buildInfo.With(prometheus.Labels{
		"version":    version,
		"build_date": buildDate,
		"go_version": runtime.Version(),
		"env":        env,
	}).Set(1)
}
