package histogram

import (
	"fmt"
	"slices"
)

// Bucket counts a partion of values.
type Bucket struct {
	// Count is the number of values represented in the bucket.
	Count int
	// Min is the low, inclusive bound of the bucket.
	Min float64
	// Max is the high, exclusive bound of the bucket.
	Max float64
}

// BSArbitrarySpan returns Buckets conforming to start, stop, and nbins.
func BSArbitrarySpan(start float64, stop float64, nbins int) func([]float64) ([]Bucket, error) {
	return func(input []float64) ([]Bucket, error) {
		return Linspace(start, stop, nbins)
	}
}

// BSExactSpan returns Buckets that span the input array from the smallest input
// to the largest input exactly.
func BSExactSpan(nb int) func([]float64) ([]Bucket, error) {
	return func(input []float64) ([]Bucket, error) {
		min := slices.Min(input)
		max := slices.Max(input)
		return Linspace(min, max, nb)
	}
}

// Bucketer is a function that takes a value and returns the bucket index in the
// supplied slice of Buckets. Implementations should return -1 if the value does
// not fit in a bucket. Callers MUST ensure that bs is sorted.
type Bucketer func(val float64, bs []Bucket) int

// DefaultBucketer uses BinearySearchFunc to bucket the supplied value.
func DefaultBucketer(val float64, bs []Bucket) int {
	idx, found := slices.BinarySearchFunc(bs, val, func(b Bucket, v float64) int {
		// return found if value falls in bucket OR if the value falls on the
		// edge of the final bucket
		if (b.Min <= v && b.Max > v) || (b == bs[len(bs)-1] && v == b.Max) {
			return 0
		}
		if b.Max <= v {
			return -1
		}
		if b.Min > v {
			return 1
		}
		// fallthrough will panic since this is never supposed to happen
		panic(fmt.Sprintf("unexpected value fallthrough: (val: %f, bs: %v)", val, bs))
	})
	if !found {
		return -1
	}
	return idx
}
