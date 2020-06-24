package main

import (
	"encoding/json"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

func (e *Exporter) collectScanMetric(ch chan<- prometheus.Metric) bool {

	type scanMetric struct {
		Total     float64
		Completed float64
		metrics   []interface{}
		Requester string
		Ongoing   bool
	}

        var ep string
        if e.opts.apiversion == "1" {
          ep = "/api/scans/all/metrics"
        } else {
          ep = "/api/v2.0/scans/all/metrics"
        }

        body := e.client.request(ep)

	var data scanMetric

	if err := json.Unmarshal(body, &data); err != nil {
		level.Error(e.logger).Log(err.Error())
		return false
	}

	scan_requester, _ := strconv.ParseFloat(data.Requester, 64)
	ch <- prometheus.MustNewConstMetric(
		scanRequesterCount, prometheus.GaugeValue, float64(scan_requester),
	)

	ch <- prometheus.MustNewConstMetric(
		scanTotalCount, prometheus.GaugeValue, float64(data.Total),
	)

	ch <- prometheus.MustNewConstMetric(
		scanCompletedCount, prometheus.GaugeValue, float64(data.Completed),
	)
	return true
}
