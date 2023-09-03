package bench

import (
	"fmt"
	"testing"
)

func BenchmarkTestA(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str += "a"
	}
}

func BenchmarkTestB(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str = fmt.Sprintf("%s%s", str, "a")
	}
}
