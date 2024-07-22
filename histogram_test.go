package histogram

import (
	"os"
	"testing"
	"time"

	"github.com/matryer/is"
)

// Tests that the arbitrary binning source works as expected using 10 equal
// sized bins from [0, 1]. For all bins (except for the final bin), the bin
// intervals are [start, stop); the last bin is [start, stop]. This replicates
// the behavior of numpy.histogram(data, range=(0, 1), bins=10).
func TestHistogramFixedBins(t *testing.T) {
	is := is.New(t)
	data := []float64{
		0.1,
		0.2, 0.21, 0.22, 0.22,
		0.3,
		0.4,
		0.5, 0.51, 0.52, 0.53, 0.54, 0.55, 0.56, 0.57, 0.58,
		0.6,
		0.8,
		0.9,
		1.0,
	}
	bs, err := BSArbitrarySpan(0, 1, 10)(data)
	is.NoErr(err)

	hist, err := Hist(data, bs, DefaultBucketer)
	is.NoErr(err)
	Fprint(os.Stdout, hist, Linear(5))

	expected := []struct {
		bl float64
		bu float64
		c  int
	}{
		{0.0, 0.1, 0},
		{0.1, 0.2, 1},
		{0.2, 0.3, 5},
		{0.3, 0.4, 0},
		{0.4, 0.5, 1},
		{0.5, 0.6, 10},
		{0.6, 0.7, 0},
		{0.7, 0.8, 0},
		{0.8, 0.9, 1},
		{0.9, 1.0, 2},
	}
	for i, b := range hist.Buckets {
		is.True(almostEqual(expected[i].bl, b.Min, 1.e-5))
		is.True(almostEqual(expected[i].bu, b.Max, 1.e-5))
		is.Equal(expected[i].c, b.Count)
	}

}

// These tests verify the histogramming behavior is consistent with the package
// github.com/aybabtme/uniplot (when configured appropriately).
func TestHistogramUniplotEquivalent(t *testing.T) {
	is := is.New(t)

	data := []float64{
		0.1,
		0.2, 0.21, 0.22, 0.22,
		0.3,
		0.4,
		0.5, 0.51, 0.52, 0.53, 0.54, 0.55, 0.56, 0.57, 0.58,
		0.6,
		0.8,
		0.9,
		1.0,
	}
	bs, err := BSExactSpan(9)(data)
	is.NoErr(err)

	hist, err := Hist(data, bs, DefaultBucketer)
	is.NoErr(err)

	Fprint(os.Stdout, hist, Linear(5))
	// 0.1-0.2  5%   ▋       1
	// 0.2-0.3  25%  ██▊     5
	// 0.3-0.4  0%   ▏
	// 0.4-0.5  5%   ▋       1
	// 0.5-0.6  45%  █████▏  9
	// 0.6-0.7  5%   ▋       1
	// 0.7-0.8  0%   ▏
	// 0.8-0.9  5%   ▋       1
	// 0.9-1    10%  █▏      2

	expected := []struct {
		bl float64
		bu float64
		c  int
	}{
		{0.1, 0.2, 1},
		{0.2, 0.3, 5},
		{0.3, 0.4, 0},
		{0.4, 0.5, 1},
		{0.5, 0.6, 9},
		{0.6, 0.7, 1},
		{0.7, 0.8, 0},
		{0.8, 0.9, 1},
		{0.9, 1.0, 2},
	}
	for i, b := range hist.Buckets {
		is.True(almostEqual(expected[i].bl, b.Min, 1.e-5))
		is.True(almostEqual(expected[i].bu, b.Max, 1.e-5))
		is.Equal(expected[i].c, b.Count)
	}

	data = []float64{
		float64(time.Millisecond * 100),
		float64(time.Millisecond * 200),
		float64(time.Millisecond * 210),
		float64(time.Millisecond * 220),
		float64(time.Millisecond * 221),
		float64(time.Millisecond * 222),
		float64(time.Millisecond * 223),
		float64(time.Millisecond * 300),
		float64(time.Millisecond * 400),
		float64(time.Millisecond * 500),
		float64(time.Millisecond * 510),
		float64(time.Millisecond * 520),
		float64(time.Millisecond * 530),
		float64(time.Millisecond * 540),
		float64(time.Millisecond * 550),
		float64(time.Millisecond * 560),
		float64(time.Millisecond * 570),
		float64(time.Millisecond * 580),
		float64(time.Millisecond * 600),
		float64(time.Millisecond * 800),
		float64(time.Millisecond * 900),
		float64(time.Millisecond * 1000),
	}
	bs, err = BSExactSpan(9)(data)
	is.NoErr(err)

	hist, err = Hist(data, bs, DefaultBucketer)
	is.NoErr(err)

	err = Fprintf(os.Stdout, hist, Linear(5), func(v float64) string {
		return time.Duration(v).String()
	})
	is.NoErr(err)
	// 100ms-200ms  4.55%  ▋       1
	// 200ms-300ms  27.3%  ███▍    6
	// 300ms-400ms  4.55%  ▋       1
	// 400ms-500ms  4.55%  ▋       1
	// 500ms-600ms  40.9%  █████▏  9
	// 600ms-700ms  4.55%  ▋       1
	// 700ms-800ms  0%     ▏
	// 800ms-900ms  4.55%  ▋       1
	// 900ms-1s     9.09%  █▏      2

	expected = []struct {
		bl float64
		bu float64
		c  int
	}{
		{float64(100 * time.Millisecond), float64(200 * time.Millisecond), 1},
		{float64(200 * time.Millisecond), float64(300 * time.Millisecond), 6},
		{float64(300 * time.Millisecond), float64(400 * time.Millisecond), 1},
		{float64(400 * time.Millisecond), float64(500 * time.Millisecond), 1},
		{float64(500 * time.Millisecond), float64(600 * time.Millisecond), 9},
		{float64(600 * time.Millisecond), float64(700 * time.Millisecond), 1},
		{float64(700 * time.Millisecond), float64(800 * time.Millisecond), 0},
		{float64(800 * time.Millisecond), float64(900 * time.Millisecond), 1},
		{float64(900 * time.Millisecond), float64(1000 * time.Millisecond), 2},
	}
	for i, b := range hist.Buckets {
		is.True(almostEqual(expected[i].bl, b.Min, 1.e-5))
		is.True(almostEqual(expected[i].bu, b.Max, 1.e-5))
		is.Equal(expected[i].c, b.Count)
	}
}

func TestBinSources(t *testing.T) {
	is := is.New(t)
	vs := []float64{1, 2, 3}

	// no error for good input
	_, err := BSExactSpan(10)(vs)
	is.NoErr(err)

	// no error for good input
	_, err = BSArbitrarySpan(0, 1, 10)(vs)
	is.NoErr(err)

	// err for bad bin count
	_, err = BSExactSpan(-1)(vs)
	is.True(err != nil)
	_, err = BSExactSpan(0)(vs)
	is.True(err != nil)
	_, err = BSArbitrarySpan(0, 1, -1)(vs)
	is.True(err != nil)
	_, err = BSArbitrarySpan(0, 1, 0)(vs)
	is.True(err != nil)

	// err for start > stop
	_, err = BSArbitrarySpan(1, 0, 10)(vs)
	is.True(err != nil)
	_, err = BSArbitrarySpan(1, 1, 10)(vs)
	is.True(err != nil)

}
