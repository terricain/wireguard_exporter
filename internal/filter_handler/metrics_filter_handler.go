package filter_handler

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/terrycain/wireguard_exporter/internal/wireguard_collector"
	"net/http"
)

type HandlerContext struct {
	MetricsPath string
	DisableExporterMetrics bool
	MaxRequests int
	Logger log.Logger
	exporterMetricsRegistry *prometheus.Registry
	MetricsHandler http.Handler
	FriendlyNames map[string]string
}

func (c *HandlerContext) IndexFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html>
			<head><title>Wireguard Exporter</title></head>
			<body>
			<h1>Wireguard Exporter</h1>
			<p><a href="` + c.MetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
}

func (c *HandlerContext) Load() {
	c.exporterMetricsRegistry = prometheus.NewRegistry()
	if !c.DisableExporterMetrics {
		c.exporterMetricsRegistry.MustRegister(
			prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
			prometheus.NewGoCollector(),
		)
	}

	if innerHandler, err := c.innerHandler(); err != nil {
		panic(fmt.Sprintf("Couldn't create metrics handler: %s", err))
	} else {
		c.MetricsHandler = innerHandler
	}
}


// Looted from node_exporter
func (c *HandlerContext) innerHandler() (http.Handler, error) {
	wg, err := wireguard_collector.NewWireguardCollector(c.Logger, c.FriendlyNames)
	if err != nil {
		return nil, fmt.Errorf("couldn't create collector: %s", err)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(version.NewCollector("wireguard"))
	if err := r.Register(wg); err != nil {
		return nil, fmt.Errorf("couldn't register wireguard collector: %s", err)
	}
	handler := promhttp.HandlerFor(
		prometheus.Gatherers{c.exporterMetricsRegistry, r},
		promhttp.HandlerOpts{
			ErrorHandling:       promhttp.ContinueOnError,
			MaxRequestsInFlight: c.MaxRequests,
			Registry:            c.exporterMetricsRegistry,
		},
	)
	if !c.DisableExporterMetrics {
		// Note that we have to use h.exporterMetricsRegistry here to
		// use the same promhttp metrics for all expositions.
		handler = promhttp.InstrumentMetricHandler(
			c.exporterMetricsRegistry, handler,
		)
	}
	return handler, nil
}