package histogram

import (
	"fmt"
	"math"
)

func Linspace(start float64, stop float64, nb int) ([]Bucket, error) {
	if nb < 1 {
		return nil, fmt.Errorf("bad value for number of bins: %d", nb)
	}
	if start >= stop {
		return nil, fmt.Errorf("bad value for start/stop: (%f, %f)", start, stop)
	}
	scale := (stop - start) / float64(nb)
	buckets := make([]Bucket, nb)
	for i := range buckets {
		bmin, bmax := start+float64(i)*scale, start+float64(i+1)*scale
		buckets[i] = Bucket{Min: bmin, Max: bmax}
	}
	return buckets, nil
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}
