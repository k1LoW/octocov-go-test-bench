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
				Name:  "number of iterations",
				Key:   "N",
				Value: float64(bench.N),
			})
			if (bench.Measured & parse.NsPerOp) != 0 {
				cs.Metrics = append(cs.Metrics, &report.CustomMetric{
					Name:  "nanoseconds per iteration",
					Key:   "NsPerOp",
					Value: float64(bench.NsPerOp),
					Unit:  "ns/op",
				})
			}
			if (bench.Measured & parse.MBPerS) != 0 {
				cs.Metrics = append(cs.Metrics, &report.CustomMetric{
					Name:  "MB processed per second",
					Key:   "MBPerS",
					Value: float64(bench.MBPerS),
					Unit:  "MB/s",
				})
			}
			if (bench.Measured & parse.AllocedBytesPerOp) != 0 {
				cs.Metrics = append(cs.Metrics, &report.CustomMetric{
					Name:  "bytes allocated per iteration",
					Key:   "AllocedBytesPerOp",
					Value: float64(bench.AllocedBytesPerOp),
					Unit:  "B/op",
				})
			}
			if (bench.Measured & parse.AllocsPerOp) != 0 {
				cs.Metrics = append(cs.Metrics, &report.CustomMetric{
					Name:  "allocs per iteration",
					Key:   "AllocsPerOp",
					Value: float64(bench.AllocsPerOp),
					Unit:  "allocs/op",
				})
			}
			cset = append(cset, cs)
		}
	}
	return cset
}
