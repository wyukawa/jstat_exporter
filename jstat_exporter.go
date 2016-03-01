package main

import (
	"flag"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/log"
)

const (
	namespace = "jstat"
)

var (
	listenAddress = flag.String("web.listen-address", ":9010", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	jstatPath     = flag.String("jstat.path", "/usr/bin/jstat", "jstat path")
	targetPid     = flag.String("target.pid", "0", "target pid")
)

type Exporter struct {
	jstatPath string
	targetPid string
	newMax  prometheus.Gauge
	newCommit  prometheus.Gauge
	oldMax  prometheus.Gauge
	oldCommit  prometheus.Gauge
	metaMax  prometheus.Gauge
	metaCommit  prometheus.Gauge
	metaUsed  prometheus.Gauge
	oldUsed   prometheus.Gauge
}

func NewExporter(jstatPath string, targetPid string) *Exporter {
	return &Exporter{
		jstatPath: jstatPath,
		targetPid: targetPid,
		newMax: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "newMax",
			Help:      "newMax",
		}),
		newCommit: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "newCommit",
			Help:      "newCommit",
		}),
		oldMax: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "oldMax",
			Help:      "oldMax",
		}),
		oldCommit: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "oldCommit",
			Help:      "oldCommit",
		}),
		metaMax: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "metaMax",
			Help:      "metaMax",
		}),
		metaCommit: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "metaCommit",
			Help:      "metaCommit",
		}),
		metaUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "metaUsed",
			Help:      "metaUsed",
		}),
		oldUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "oldUsed",
			Help:      "oldUsed",
		}),
	}
}

// Describe implements the prometheus.Collector interface.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	e.newMax.Describe(ch)
	e.newCommit.Describe(ch)
	e.oldMax.Describe(ch)
	e.oldCommit.Describe(ch)
	e.metaMax.Describe(ch)
	e.metaCommit.Describe(ch)
	e.metaUsed.Describe(ch)
	e.oldUsed.Describe(ch)
}

// Collect implements the prometheus.Collector interface.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.JstatGccapacity(ch)
	e.JstatGcold(ch)
}

func (e *Exporter) JstatGccapacity(ch chan<- prometheus.Metric) {

	out, err := exec.Command(e.jstatPath, "-gccapacity", e.targetPid).Output()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i == 1 {
			parts := strings.Fields(line)
			newMax, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.newMax.Set(newMax)
			e.newMax.Collect(ch)
			newCommit, err := strconv.ParseFloat(parts[2], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.newCommit.Set(newCommit)
			e.newCommit.Collect(ch)
		}
	}
}

func (e *Exporter) JstatGcold(ch chan<- prometheus.Metric) {

	out, err := exec.Command(e.jstatPath, "-gcold", e.targetPid).Output()
	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i == 1 {
			parts := strings.Fields(line)
			metaUsed, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.metaUsed.Set(metaUsed) // MU: Metaspace utilization (kB).
			e.metaUsed.Collect(ch)
			oldUsed, err := strconv.ParseFloat(parts[5], 64)
			if err != nil {
				log.Fatal(err)
			}
			e.oldUsed.Set(oldUsed) // OU: Old space utilization (kB).
			e.oldUsed.Collect(ch)
		}
	}
}

func main() {
	flag.Parse()

	exporter := NewExporter(*jstatPath, *targetPid)
	prometheus.MustRegister(exporter)

	log.Printf("Starting Server: %s", *listenAddress)
	http.Handle(*metricsPath, prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
		<head><title>jstat Exporter</title></head>
		<body>
		<h1>jstat Exporter</h1>
		<p><a href="` + *metricsPath + `">Metrics</a></p>
		</body>
		</html>`))
	})
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}

}
