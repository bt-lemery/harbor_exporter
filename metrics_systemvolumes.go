package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) collectSystemVolumesMetric(ch chan<- prometheus.Metric) bool {
	type systemVolumesMetric struct {
		Storage struct {
			Total float64
			Free  float64
		}
	}

        var ep string
        if e.opts.apiversion == "1" {
          ep = "/api/systeminfo/volumes"
        } else {
          ep = "/api/v2.0/systeminfo/volumes"
        }

        body := e.client.request(ep)

	var data systemVolumesMetric
	if err := json.Unmarshal(body, &data); err != nil {
		level.Error(e.logger).Log(err.Error())
		return false
	}

	ch <- prometheus.MustNewConstMetric(
		systemVolumes, prometheus.GaugeValue, data.Storage.Total, "total",
	)
	ch <- prometheus.MustNewConstMetric(
		systemVolumes, prometheus.GaugeValue, data.Storage.Free, "free",
	)

	return true
}
