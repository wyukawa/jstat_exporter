package main

import (
	"flag"
	"fmt"
	"os/exec"
	"strings"

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

func main() {
	flag.Parse()

	out, err := exec.Command(*jstatPath, "-gcold", *targetPid).Output()

	if err != nil {
		log.Fatal(err)
	}

	for i, line := range strings.Split(string(out), "\n") {
		if i == 1 {
			parts := strings.Fields(line)
			fmt.Println(parts[1]) // MU: Metaspace utilization (kB).
			fmt.Println(parts[5]) // OU: Old space utilization (kB).
		}
	}

}
