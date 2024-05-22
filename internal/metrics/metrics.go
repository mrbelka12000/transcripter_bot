package metrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ExposeMetrics(port int) error {
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
