package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const namespace = "fluxcd"

var (
	listenAddress = flag.String("web.listen-address", ":8080",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")

	// Metrics
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last fluxcd query successful",
		nil, nil,
	)
	versionsApplied = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "version_applied"),
		"List the current installed apps version",
		[]string{"namespace", "name", "ready", "message", "revision", "suspended"}, nil,
	)
)

type Exporter struct{}

func NewExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- versionsApplied
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	out, err := exec.Command("flux", "get", "-A", "all", "--status-selector", "ready=true", "--no-header").Output()

	if err != nil {
		ch <- prometheus.MustNewConstMetric(
			up, prometheus.GaugeValue, 0,
		)
		log.Fatal(err)
		return
	}
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)

	results := strings.Split(string(out), "\n")

	for i, s := range results {
		if len(s) > 0 {
			fields := strings.Fields(s)
			fmt.Println(i, fields)
			ch <- prometheus.MustNewConstMetric(
				versionsApplied, prometheus.GaugeValue, 1, fields[0], fields[1], fields[2], fields[3], fields[6], fields[7],
			)
		}
	}

	fmt.Println("Metrics collected")
}

func main() {
	log.SetOutput(os.Stdout)
	flag.Parse()

	exporter := NewExporter()
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Fluxcd Version Exporter</title></head>
             <body>
             <h1>Fluxcd Version Exporter</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
	fmt.Println("Exporter started")
}
