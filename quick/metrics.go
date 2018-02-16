package main

import (
	"github.com/ohsu-comp-bio/mortar/graph"
	"github.com/prometheus/client_golang/prometheus"
)

var completeGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "ohsu",
		Subsystem: "mortar",
		Name:      "steps_complete",
		Help:      "Number of steps complete per workflow run.",
	},
	[]string{"run"},
)

var totalGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Namespace: "ohsu",
		Subsystem: "mortar",
		Name:      "total_steps",
		Help:      "Total number of steps in a run.",
	},
	[]string{"run"},
)

func init() {
	prometheus.MustRegister(completeGauge)
	prometheus.MustRegister(totalGauge)
}

func updateMetrics(cli graph.Client, graphID string) {
	d := getData(cli, graphID)
	for rid, run := range d.Runs {
		completeGauge.WithLabelValues(rid).Set(float64(run.Complete))
		totalGauge.WithLabelValues(rid).Set(float64(run.Total))
	}
}
