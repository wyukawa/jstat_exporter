# jstat Exporter for Prometheus
Exports jstat result for Prometheus consumption.

How to build
```
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/log
go build jstat_exporter.go
```

Help on flags of jstat_exporter:
```
-web.listen-address string
    Address on which to expose metrics and web interface. (default ":9010")
-web.telemetry-path string
    Path under which to expose metrics. (default "/metrics")
```

Tested on JDK8
