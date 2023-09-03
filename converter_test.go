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
		in   parse.Set
		want []*report.CustomMetricSet
	}{
		{parse.Set{}, nil},
		{
			parse.Set{
				"Benchmark-0": []*parse.Benchmark{
					{
						Name:              "Benchmark-0",
						N:                 1,
						NsPerOp:           100,
						AllocedBytesPerOp: 25,
						AllocsPerOp:       50,
					},
				},
			},
			[]*report.CustomMetricSet{
				{
					Name: "Benchmark-0",
					Key:  "Benchmark-0",
					Metrics: []*report.CustomMetric{
						{Name: "N", Key: "N", Value: 1},
						{Name: "ns/op", Key: "NsPerOp", Value: 100, Unit: "ns/op"},
						{Name: "B/op", Key: "AllocedBytesPerOp", Value: 25, Unit: "B/op"},
						{Name: "allocs/op", Key: "AllocsPerOp", Value: 50, Unit: "allocs/op"},
					},
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got := Convert(tt.in)
			opts := []cmp.Option{
				cmpopts.IgnoreFields(report.CustomMetricSet{}, "report"),
			}
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}
