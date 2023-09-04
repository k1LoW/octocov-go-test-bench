package gotestbench

import (
	"fmt"
	"sort"

	"github.com/k1LoW/octocov/report"
	"github.com/samber/lo"
	"golang.org/x/tools/benchmark/parse"
)

type benchGroup struct {
	key    string
	benchs []*parse.Benchmark
	n      int
}

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
		var names []string
		g := map[string]*benchGroup{}
		for _, b := range benchs {
			names = append(names, b.Name)
			bg, ok := g[b.Name]
			if !ok {
				g[b.Name] = &benchGroup{
					key:    b.Name,
					benchs: []*parse.Benchmark{b},
					n:      1,
				}
				continue
			}
			bg.benchs = append(bg.benchs, b)
			bg.n++
		}
		names = unique(names)
		for _, n := range names {
			bg := g[n]
			name := n
			if bg.n > 1 {
				name = fmt.Sprintf("%s (average of %d)", name, bg.n)
			}
			cs := &report.CustomMetricSet{
				Name: name,
				Key:  n,
			}
			b := bg.benchs[0]
			cs.Metrics = append(cs.Metrics, &report.CustomMetric{
				Name: "Number of iterations",
				Key:  "N",
				Value: lo.Reduce(bg.benchs, func(agg float64, item *parse.Benchmark, _ int) float64 {
					return agg + float64(b.N)
				}, 0.0) / float64(bg.n),
			})
			if (b.Measured & parse.NsPerOp) != 0 {
				cs.Metrics = append(cs.Metrics, &report.CustomMetric{
					Name: "Nanoseconds per iteration",
					Key:  "NsPerOp",
					Value: lo.Reduce(bg.benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
						return agg + float64(b.NsPerOp)
					}, 0.0) / float64(bg.n),
					Unit: " ns/op",
				})
			}
			if (b.Measured & parse.MBPerS) != 0 {
				cs.Metrics = append(cs.Metrics, &report.CustomMetric{
					Name: "MB processed per second",
					Key:  "MBPerS",
					Value: lo.Reduce(bg.benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
						return agg + float64(b.MBPerS)
					}, 0.0) / float64(bg.n),
					Unit: " MB/s",
				})
			}
			if (b.Measured & parse.AllocedBytesPerOp) != 0 {
				cs.Metrics = append(cs.Metrics, &report.CustomMetric{
					Name: "Bytes allocated per iteration",
					Key:  "AllocedBytesPerOp",
					Value: lo.Reduce(bg.benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
						return agg + float64(b.AllocedBytesPerOp)
					}, 0.0) / float64(bg.n),
					Unit: " B/op",
				})
			}
			if (b.Measured & parse.AllocsPerOp) != 0 {
				cs.Metrics = append(cs.Metrics, &report.CustomMetric{
					Name: "Allocs per iteration",
					Key:  "AllocsPerOp",
					Value: lo.Reduce(bg.benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
						return agg + float64(b.AllocsPerOp)
					}, 0.0) / float64(bg.n),
					Unit: " allocs/op",
				})
			}
			cset = append(cset, cs)
		}
	}
	return cset
}

func unique(s []string) []string {
	m := map[string]struct{}{}
	for _, v := range s {
		m[v] = struct{}{}
	}
	var ret []string
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}
