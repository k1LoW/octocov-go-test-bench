package gotestbench

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k1LoW/octocov/report"
	"golang.org/x/tools/benchmark/parse"
)

func TestConverter(t *testing.T) {
	tests := []struct {
		in    parse.Set
		stype string
		want  []*report.CustomMetricSet
	}{
		{parse.Set{}, StasticsTypeAvg, nil},
		{
			parse.Set{
				"Benchmark-0": []*parse.Benchmark{
					{
						Name:              "Benchmark-0",
						N:                 1,
						NsPerOp:           100,
						AllocedBytesPerOp: 25,
						AllocsPerOp:       50,
						Measured:          parse.NsPerOp | parse.AllocedBytesPerOp | parse.AllocsPerOp,
					},
				},
			},
			StasticsTypeAvg,
			[]*report.CustomMetricSet{
				{
					Name: "Benchmark-0",
					Key:  "Benchmark-0",
					Metrics: []*report.CustomMetric{
						{Name: "Number of iterations", Key: "N", Value: 1},
						{Name: "Nanoseconds per iteration", Key: "NsPerOp", Value: 100, Unit: " ns/op"},
						{Name: "Bytes allocated per iteration", Key: "AllocedBytesPerOp", Value: 25, Unit: " B/op"},
						{Name: "Allocs per iteration", Key: "AllocsPerOp", Value: 50, Unit: " allocs/op"},
					},
				},
			},
		},
		{
			parse.Set{
				"Benchmark-0": []*parse.Benchmark{
					{
						Name:              "Benchmark-0",
						N:                 1,
						NsPerOp:           100,
						AllocedBytesPerOp: 25,
						AllocsPerOp:       50,
						Measured:          0,
					},
				},
			},
			StasticsTypeAvg,
			[]*report.CustomMetricSet{
				{
					Name: "Benchmark-0",
					Key:  "Benchmark-0",
					Metrics: []*report.CustomMetric{
						{Name: "Number of iterations", Key: "N", Value: 1},
					},
				},
			},
		},
		{
			parse.Set{
				"Benchmark-0": []*parse.Benchmark{
					{
						Name:              "Benchmark-0",
						N:                 1,
						NsPerOp:           100,
						AllocedBytesPerOp: 25,
						AllocsPerOp:       50,
						Measured:          parse.NsPerOp | parse.AllocedBytesPerOp | parse.AllocsPerOp,
					},
					{
						Name:              "Benchmark-0",
						N:                 1,
						NsPerOp:           10,
						AllocedBytesPerOp: 25,
						AllocsPerOp:       100,
						Measured:          parse.NsPerOp | parse.AllocedBytesPerOp | parse.AllocsPerOp,
					},
				},
			},
			StasticsTypeAvg,
			[]*report.CustomMetricSet{
				{
					Name: "Benchmark-0 (average of 2 benchmarks)",
					Key:  "Benchmark-0",
					Metrics: []*report.CustomMetric{
						{Name: "Number of iterations", Key: "N", Value: 1},
						{Name: "Nanoseconds per iteration", Key: "NsPerOp", Value: 55, Unit: " ns/op"},
						{Name: "Bytes allocated per iteration", Key: "AllocedBytesPerOp", Value: 25, Unit: " B/op"},
						{Name: "Allocs per iteration", Key: "AllocsPerOp", Value: 75, Unit: " allocs/op"},
					},
				},
			},
		},
		{
			parse.Set{
				"Benchmark-0": []*parse.Benchmark{
					{
						Name:              "Benchmark-0",
						N:                 1,
						NsPerOp:           100,
						AllocedBytesPerOp: 25,
						AllocsPerOp:       50,
						Measured:          parse.NsPerOp | parse.AllocedBytesPerOp | parse.AllocsPerOp,
					},
					{
						Name:              "Benchmark-0",
						N:                 1,
						NsPerOp:           10,
						AllocedBytesPerOp: 25,
						AllocsPerOp:       100,
						Measured:          parse.NsPerOp | parse.AllocedBytesPerOp | parse.AllocsPerOp,
					},
				},
			},
			StasticsTypeMedByNsPerOp,
			[]*report.CustomMetricSet{
				{
					Name: "Benchmark-0 (median by \"Nanoseconds per iteration\" of 2 benchmarks)",
					Key:  "Benchmark-0",
					Metrics: []*report.CustomMetric{
						{Name: "Number of iterations", Key: "N", Value: 1},
						{Name: "Nanoseconds per iteration", Key: "NsPerOp", Value: 55, Unit: " ns/op"},
						{Name: "Bytes allocated per iteration", Key: "AllocedBytesPerOp", Value: 25, Unit: " B/op"},
						{Name: "Allocs per iteration", Key: "AllocsPerOp", Value: 75, Unit: " allocs/op"},
					},
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, err := Convert(tt.in, tt.stype)
			if err != nil {
				t.Fatal(err)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(report.CustomMetricSet{}, "report"),
			}
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}
