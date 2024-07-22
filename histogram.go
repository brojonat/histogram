package histogram

// Histogram holds a count of values partionned over buckets.
type Histogram struct {
	// Min is the size of the smallest bucket.
	Min int
	// Max is the size of the biggest bucket.
	Max int
	// Count is the total size of all buckets.
	Count int
	// Buckets over which values are partionned.
	Buckets []Bucket
}

// Scale gives the scaled count of the bucket at idx, using the
// provided scale func.
func (h Histogram) Scale(s ScaleFunc, idx int) float64 {
	bkt := h.Buckets[idx]
	scale := s(h.Min, h.Max, bkt.Count)
	return scale
}

// ScaleFunc is the type to implement to scale an histogram.
type ScaleFunc func(min, max, value int) float64

// Linear builds a ScaleFunc that will linearly scale the values of
// an histogram so that they do not exceed width.
func Linear(width int) ScaleFunc {
	return func(min, max, value int) float64 {
		if min == max {
			return 1
		}
		return float64(value-min) / float64(max-min) * float64(width)
	}
}

func Hist(input []float64, buckets []Bucket, bf Bucketer) (*Histogram, error) {
	if len(buckets) == 0 {
		return &Histogram{}, nil
	}

	if len(buckets) == 1 {
		buckets[0].Count = len(input)
		return &Histogram{
			Min:     len(input),
			Max:     len(input),
			Count:   len(input),
			Buckets: buckets,
		}, nil
	}

	minC, maxC := 0, 0
	for _, val := range input {
		bi := bf(val, buckets)
		if bi < 0 {
			continue
		}
		buckets[bi].Count++
		minC = imin(minC, buckets[bi].Count)
		maxC = imax(maxC, buckets[bi].Count)
	}

	return &Histogram{
		Min:     minC,
		Max:     maxC,
		Count:   len(input),
		Buckets: buckets,
	}, nil
}
