package bloom

import (
    "fmt"
	"testing"
)

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for errorRate := 0.1; errorRate > 0.0001; errorRate /= 10 {
			New(500000, errorRate)
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	filter := New(500000, 0.001)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Add([]byte(fmt.Sprintf("T%v", i)))
	}
}

func BenchmarkContains(b *testing.B) {
	filter := New(500000, 0.001)
	for i := 0; i < b.N; i++ {
		filter.Add([]byte(fmt.Sprintf("T%v", i)))
	}
	b.ResetTimer()

	for i := b.N / 2; i < 3*b.N/2; i++ {
		filter.Contains([]byte(fmt.Sprintf("T%v", i)))
	}
}

func TestWorkAsExpected(t *testing.T) {
	capacity := 100000.0
	errorRate := 0.001

	filter := New(int(capacity), errorRate)

	falseNegative := 0.0
	falsePositive := 0.0

	limit := capacity * 3
	for k := 0.0; k < limit; k += 3.0 {
		filter.Add([]byte(fmt.Sprintf("T%v", k)))
	}

	samples := capacity * 9
	for k := 0.0; k < samples; k++ {
		if int(k)%3 == 0 && k < limit {
			falseNegative += boolToFloat(!filter.Contains([]byte(fmt.Sprintf("T%v", k))))
		} else if k > limit {
			falsePositive += boolToFloat(filter.Contains([]byte(fmt.Sprintf("T%v", k))))
		}
	}

	observedError := falsePositive / samples
	pass := observedError < errorRate && falseNegative == 0

	if !pass {
		t.Fatal("Bloom filter failing to keep zero false negatives and expected false positive ration.")
	}
}

func boolToFloat(b bool) float64 {
    if b {
        return 1.0
    }
    return 0.0
}
