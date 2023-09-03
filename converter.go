package gotestbench

import (
	"sort"

	"github.com/k1LoW/octocov/report"
	"golang.org/x/tools/benchmark/parse"
)

func Convert(set parse.Set) []*report.CustomMetricSet {
	var keys []string
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var cset []*report.CustomMetricSet
	for _, k := range keys {
		benchs, ok := set[k]
		if !ok {
			continue
		}
		for _, bench := range benchs {
			cs := &report.CustomMetricSet{
				Name: bench.Name,
				Key:  bench.Name,
			}
			cs.Metrics = append(cs.Metrics, &report.CustomMetric{
				Name:  "N",
				Key:   "N",
				Value: float64(bench.N),
			})
			cs.Metrics = append(cs.Metrics, &report.CustomMetric{
				Name:  "ns/op",
				Key:   "NsPerOp",
				Value: float64(bench.NsPerOp),
				Unit:  "ns/op",
			})
			cs.Metrics = append(cs.Metrics, &report.CustomMetric{
				Name:  "B/op",
				Key:   "AllocedBytesPerOp",
				Value: float64(bench.AllocedBytesPerOp),
				Unit:  "B/op",
			})
			cs.Metrics = append(cs.Metrics, &report.CustomMetric{
				Name:  "allocs/op",
				Key:   "AllocsPerOp",
				Value: float64(bench.AllocsPerOp),
				Unit:  "allocs/op",
			})
			cset = append(cset, cs)
		}
	}
	return cset
}
