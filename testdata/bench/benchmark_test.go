package bench

import (
	"fmt"
	"testing"
)

func BenchmarkA(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str += "a"
	}
}

func BenchmarkB(b *testing.B) {
	str := ""
	for i := 0; i < b.N; i++ {
		str = fmt.Sprintf("%s%s", str, "a")
	}
}
