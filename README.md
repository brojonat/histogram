# histogram

This package provides histogramming functionality. It is based on `https://github.com/aybabtme/uniplot`, but has some enhancements that allow callers to supply their own binning and indexing implementations.

## Examples

See the tests for more examples, but the basic usage is as follows:

```go
vs := []float64{0.1, 0.15, 0.1, 0.2, 0.7, 0.4, 0.71}
bs := BSArbitrarySpan(0, 1, 10)(vs)
h, err := Hist(vs, bs, DefaultBucketer)
if err != nil {
    // do something
}
Fprint(os.Stdout, hist, Linear(5))
```
