package histogram

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

type Indexable interface {
	constraints.Unsigned | constraints.Signed | constraints.Float
}

// Bucket counts a partion of values.
type Bucket[T Indexable] struct {
	// Min is the low, inclusive bound of the bucket.
	Min T
	// Max is the high, exclusive bound of the bucket. If
	// this bucket is the last bucket, the bound is inclusive
	// and contains the max value of the histogram.
	Max T
	// Count is the number of values represented in the bucket.
	Count int
}

func (b Bucket[T]) within(val T, inclusiveMax bool) bool {
	if inclusiveMax {
		return b.Min <= val && val <= b.Max
	}
	return b.Min <= val && val < b.Max
}

// Histogram holds a count of values partioned over buckets.
type Histogram[T Indexable] struct {
	// Min is the size of the smallest bucket.
	Min int
	// Max is the size of the biggest bucket.
	Max int
	// Count is the total size of all buckets.
	Count int
	// Buckets over which values are partionned.
	Buckets []Bucket[T]
}

// Bucket finds the index of the Bucket set that a value falls into
func (h Histogram[T]) Index(val T) (int, error) {
	last := len(h.Buckets) - 1
	for i, b := range h.Buckets {
		if b.within(val, last == i) {
			return i, nil
		}
	}
	// check last bucket edge case
	return len(h.Buckets) - 1, fmt.Errorf("value outside of histogram range")
}

func newHistogram[T Indexable](buckets []Bucket[T], data []T) Histogram[T] {
	h := Histogram[T]{
		Buckets: buckets,
	}

	last := len(buckets) - 1
	for _, val := range data {
		for i, bucket := range h.Buckets {
			if bucket.within(val, i == last) {
				h.Count++
				h.Buckets[i].Count++
				h.Min = _min(h.Min, h.Buckets[i].Count)
				h.Max = _max(h.Max, h.Buckets[i].Count)
				break
			}
		}
	}

	return h
}

// Create creates an histogram partioning input over `bins` buckets.
func Create[T Indexable](bins int, input []T) Histogram[T] {
	count := len(input)

	if count == 0 || bins == 0 {
		return Histogram[T]{}
	}

	sort.SliceStable(input, func(i, j int) bool {
		return i < j
	})

	min := input[0]
	max := input[len(input)-1]

	scale := float64(max-min) / float64(bins)

	if min == max {
		bins = 1
	}

	buckets := make([]Bucket[T], bins)
	for i := 0; i < bins; i++ {
		buckets[i] = Bucket[T]{
			Min: T(float64(i)*scale) + min,
			Max: T(float64(i+1)*scale) + min,
		}
	}

	return newHistogram(buckets, input)
}

// CreateRanged creates an histogram by custom range, like elasticsearch data histogram query(left closed)
func CreateRanged[T Indexable](min, max, interval T, input []T) Histogram[T] {
	if interval == 0 {
		return Histogram[T]{}
	}

	sort.SliceStable(input, func(i, j int) bool {
		return i < j
	})

	bins := 1
	if max != min {
		bins = int(float64(max-min) / float64(interval))
	}

	buckets := make([]Bucket[T], bins)

	for i := 0; i < bins; i++ {
		buckets[i] = Bucket[T]{
			Min: min + (interval * T(i)),
			Max: _min(min+(interval*T(i+1)), max),
		}
	}

	return newHistogram(buckets, input)
}

// CreateLog creates an histogram with logarithmic partioning
func CreateLog[T Indexable](power float64, input []T) Histogram[T] {
	count := len(input)
	if count == 0 || power <= 0 {
		return Histogram[T]{}
	}

	sort.SliceStable(input, func(i, j int) bool {
		return i < j
	})

	minValue := input[0]
	maxValue := input[len(input)-1]

	fromPower := math.Floor(logbase(minValue, power))
	toPower := math.Floor(logbase(maxValue, power))

	buckets := make([]Bucket[T], int(toPower-fromPower)+1)
	for i := 0; i < len(buckets); i++ {
		buckets[i].Min = T(math.Pow(power, float64(i)+fromPower))
		buckets[i].Max = T(math.Pow(power, float64(i+1)+fromPower))
	}

	return newHistogram(buckets, input)
}

func logbase[T Indexable](a T, base float64) float64 {
	return math.Log2(float64(a)) / math.Log2(base)
}

func _min[T Indexable](a T, b T) T {
	if a < b {
		return a
	}
	return b
}

func _max[T Indexable](a T, b T) T {
	if a < b {
		return b
	}
	return b
}
