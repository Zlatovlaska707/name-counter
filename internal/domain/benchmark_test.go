package domain

import (
	"strings"
	"testing"
)

func BenchmarkCounter_Count(b *testing.B) {
	data := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = "Name" + string(rune(i%100))
	}
	input := strings.Join(data, "\n")

	counter := NewCounter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)
		_, err := counter.Count(r, SortByFrequencyDesc)
		if err != nil {
			b.Fatalf("Count() failed: %v", err)
		}
	}
}

func BenchmarkCounter_Count_SortByName(b *testing.B) {
	data := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		data[i] = "Name" + string(rune(i%100))
	}
	input := strings.Join(data, "\n")

	counter := NewCounter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := strings.NewReader(input)
		_, err := counter.Count(r, SortByNameAsc)
		if err != nil {
			b.Fatalf("Count() failed: %v", err)
		}
	}
}
