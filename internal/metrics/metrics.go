package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// SearchQueue ...
	SearchQueue prometheus.Gauge
	// SearchDuration ...
	SearchDuration prometheus.Histogram
	// SearchCount ...
	SearchCount prometheus.Counter
)

// Setup creates metrics ready for use
func Setup() {
	SearchQueue = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "wpdir",
		Name:      "search_queue",
		Help:      "Number of searches waiting to be processed.",
	})
	prometheus.MustRegister(SearchQueue)

	SearchDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "wpdir",
		Name:      "search_duration_seconds",
		Help:      "Time taken to complete searches",
	})
	prometheus.MustRegister(SearchDuration)

	SearchCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "wpdir",
		Name:      "search_count",
		Help:      "Total number of searches",
	})
	prometheus.MustRegister(SearchCount)
}
