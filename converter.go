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
}

const (
	StasticsTypeAvg                    = "avg"
	StasticsTypeMedByN                 = "med-by-n"
	StasticsTypeMedByNsPerOp           = "med-by-ns-per-op"
	StasticsTypeMedByMBPerS            = "med-by-mb-per-s"
	StasticsTypeMedByAllocedBytesPerOp = "med-by-alloced-bytes-per-op"
	StasticsTypeMedByAllocsPerOp       = "med-by-allocs-per-op"
)

func Convert(set parse.Set, stype string) ([]*report.CustomMetricSet, error) {
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
				}
				continue
			}
			bg.benchs = append(bg.benchs, b)
		}
		names = unique(names)
		for _, n := range names {
			bg := g[n]
			name := n
			var (
				metrics []*report.CustomMetric
				err     error
			)
			switch stype {
			case StasticsTypeAvg:
				if len(bg.benchs) > 1 {
					name = fmt.Sprintf("%s (average of %d benchmarks)", name, len(bg.benchs))
				}
				metrics, err = avg(bg.benchs)
				if err != nil {
					return nil, err
				}
			case StasticsTypeMedByN:
				if len(bg.benchs) > 1 {
					name = fmt.Sprintf("%s (median by %q of %d benchmarks)", name, "Number of iterations", len(bg.benchs))
				}
				metrics, err = medByN(bg.benchs)
				if err != nil {
					return nil, err
				}
			case StasticsTypeMedByNsPerOp:
				if len(bg.benchs) > 1 {
					name = fmt.Sprintf("%s (median by %q of %d benchmarks)", name, "Nanoseconds per iteration", len(bg.benchs))
				}
				metrics, err = medByNsPerOp(bg.benchs)
				if err != nil {
					return nil, err
				}
			case StasticsTypeMedByMBPerS:
				if len(bg.benchs) > 1 {
					name = fmt.Sprintf("%s (median by %q of %d benchmarks)", name, "MB processed per second", len(bg.benchs))
				}
				metrics, err = medByMBPerS(bg.benchs)
				if err != nil {
					return nil, err
				}
			case StasticsTypeMedByAllocedBytesPerOp:
				if len(bg.benchs) > 1 {
					name = fmt.Sprintf("%s (median by %q of %d benchmarks)", name, "Bytes allocated per iteration", len(bg.benchs))
				}
				metrics, err = medByAllocedBytesPerOp(bg.benchs)
				if err != nil {
					return nil, err
				}
			case StasticsTypeMedByAllocsPerOp:
				if len(bg.benchs) > 1 {
					name = fmt.Sprintf("%s (median by %q of %d benchmarks)", name, "Allocs per iteration", len(bg.benchs))
				}
				metrics, err = medByAllocsPerOp(bg.benchs)
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unknown stastics type: %s", stype)
			}
			cs := &report.CustomMetricSet{
				Name:    name,
				Key:     n,
				Metrics: metrics,
			}
			cset = append(cset, cs)
		}
	}
	return cset, nil
}

func avg(benchs []*parse.Benchmark) ([]*report.CustomMetric, error) {
	metrics := []*report.CustomMetric{}
	b := benchs[0]
	n := len(benchs)
	metrics = append(metrics, &report.CustomMetric{
		Name: "Number of iterations",
		Key:  "N",
		Value: lo.Reduce(benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
			return agg + float64(b.N)
		}, 0.0) / float64(n),
	})
	if (b.Measured & parse.NsPerOp) != 0 {
		metrics = append(metrics, &report.CustomMetric{
			Name: "Nanoseconds per iteration",
			Key:  "NsPerOp",
			Value: lo.Reduce(benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
				return agg + float64(b.NsPerOp)
			}, 0.0) / float64(n),
			Unit: " ns/op",
		})
	}
	if (b.Measured & parse.MBPerS) != 0 {
		metrics = append(metrics, &report.CustomMetric{
			Name: "MB processed per second",
			Key:  "MBPerS",
			Value: lo.Reduce(benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
				return agg + float64(b.MBPerS)
			}, 0.0) / float64(n),
			Unit: " MB/s",
		})
	}
	if (b.Measured & parse.AllocedBytesPerOp) != 0 {
		metrics = append(metrics, &report.CustomMetric{
			Name: "Bytes allocated per iteration",
			Key:  "AllocedBytesPerOp",
			Value: lo.Reduce(benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
				return agg + float64(b.AllocedBytesPerOp)
			}, 0.0) / float64(n),
			Unit: " B/op",
		})
	}
	if (b.Measured & parse.AllocsPerOp) != 0 {
		metrics = append(metrics, &report.CustomMetric{
			Name: "Allocs per iteration",
			Key:  "AllocsPerOp",
			Value: lo.Reduce(benchs, func(agg float64, b *parse.Benchmark, _ int) float64 {
				return agg + float64(b.AllocsPerOp)
			}, 0.0) / float64(n),
			Unit: " allocs/op",
		})
	}
	return metrics, nil
}

func medByN(benchs []*parse.Benchmark) ([]*report.CustomMetric, error) {
	sort.Slice(benchs, func(i, j int) bool {
		return benchs[i].N < benchs[j].N
	})
	if len(benchs)%2 == 0 {
		return avg([]*parse.Benchmark{benchs[len(benchs)/2-1], benchs[len(benchs)/2]})
	}
	return avg([]*parse.Benchmark{benchs[len(benchs)/2]})
}

func medByNsPerOp(benchs []*parse.Benchmark) ([]*report.CustomMetric, error) {
	b := benchs[0]
	if (b.Measured & parse.NsPerOp) == 0 {
		return nil, fmt.Errorf("benchmarks do not have NsPerOp")
	}
	sort.Slice(benchs, func(i, j int) bool {
		return benchs[i].NsPerOp < benchs[j].NsPerOp
	})
	if len(benchs)%2 == 0 {
		return avg([]*parse.Benchmark{benchs[len(benchs)/2-1], benchs[len(benchs)/2]})
	}
	return avg([]*parse.Benchmark{benchs[len(benchs)/2]})
}

func medByMBPerS(benchs []*parse.Benchmark) ([]*report.CustomMetric, error) {
	b := benchs[0]
	if (b.Measured & parse.MBPerS) == 0 {
		return nil, fmt.Errorf("benchmarks do not have MBPerS")
	}
	sort.Slice(benchs, func(i, j int) bool {
		return benchs[i].MBPerS < benchs[j].MBPerS
	})
	if len(benchs)%2 == 0 {
		return avg([]*parse.Benchmark{benchs[len(benchs)/2-1], benchs[len(benchs)/2]})
	}
	return avg([]*parse.Benchmark{benchs[len(benchs)/2]})
}

func medByAllocedBytesPerOp(benchs []*parse.Benchmark) ([]*report.CustomMetric, error) {
	b := benchs[0]
	if (b.Measured & parse.AllocedBytesPerOp) == 0 {
		return nil, fmt.Errorf("benchmarks do not have AllocedBytesPerOp")
	}
	sort.Slice(benchs, func(i, j int) bool {
		return benchs[i].AllocedBytesPerOp < benchs[j].AllocedBytesPerOp
	})
	if len(benchs)%2 == 0 {
		return avg([]*parse.Benchmark{benchs[len(benchs)/2-1], benchs[len(benchs)/2]})
	}
	return avg([]*parse.Benchmark{benchs[len(benchs)/2]})
}

func medByAllocsPerOp(benchs []*parse.Benchmark) ([]*report.CustomMetric, error) {
	b := benchs[0]
	if (b.Measured & parse.AllocsPerOp) == 0 {
		return nil, fmt.Errorf("benchmarks do not have AllocsPerOp")
	}
	sort.Slice(benchs, func(i, j int) bool {
		return benchs[i].AllocsPerOp < benchs[j].AllocsPerOp
	})
	if len(benchs)%2 == 0 {
		return avg([]*parse.Benchmark{benchs[len(benchs)/2-1], benchs[len(benchs)/2]})
	}
	return avg([]*parse.Benchmark{benchs[len(benchs)/2]})
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
