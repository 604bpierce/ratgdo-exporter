package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/604pierce/ratgdo-exporter/internal/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	defaultHost := os.Getenv("HOST")
	if defaultHost == "" {
		defaultHost = "localhost"
	}

	defaultPort := os.Getenv("METRICS_PORT")
	if defaultPort == "" {
		defaultPort = "9100"
	}

	host := flag.String("host", defaultHost, "URL with schema to ratgdo host, eg: http://ratgo.local")
	port := flag.String("metrics-port", defaultPort, "Port to expose the metrics page. default is 9100")
	flag.Parse()

	reg := prometheus.NewRegistry()

	collector := collector.NewCollector(*host)
	reg.MustRegister(collector)

	// Expose /metrics
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	log.Printf("Exporter is listening on :%s/metrics\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}
