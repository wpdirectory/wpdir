package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

func metricsMiddleware(h http.Handler) http.Handler {
	fn := prometheus.InstrumentHandler(
		"middleware", h,
	)

	return http.HandlerFunc(fn)
}
