package histogram

import (
	"fmt"
	"math"
	"sort"

	"golang.org/x/exp/constraints"
)

// Indexable defines the constraints for Histogram values
type Indexable interface {
	constraints.Unsigned | constraints.Signed | constraints.Float
}

// Bucket counts a partion of values.
type Bucket[T Indexable] struct {
	inclusiveMax bool
	// Min is the low, inclusive bound of the bucket.
	Min T
	// Max is the high, exclusive bound of the bucket. If
	// this bucket is the last bucket, the bound is inclusive
	// and contains the max value of the histogram.
	Max T
	// Count is the number of values represented in the bucket.
	Count int
}

func (b Bucket[T]) within(val T) bool {
	if b.inclusiveMax {
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

// Index finds the index of the Bucket set that a value falls into
func (h Histogram[T]) Index(val T) (int, error) {
	for i, b := range h.Buckets {
		if b.within(val) {
			return i, nil
		}
	}
	// check last bucket edge case
	return len(h.Buckets) - 1, fmt.Errorf("value outside of histogram range")
}

func newHistogram[T Indexable](buckets []Bucket[T], data []T) Histogram[T] {
	// mark the final bucket inclusive
	buckets[len(buckets)-1].inclusiveMax = true
	h := Histogram[T]{
		Buckets: buckets,
	}

	for _, val := range data {
		for i, bucket := range h.Buckets {
			if bucket.within(val) {
				h.Count++
				h.Buckets[i].Count++
				h.Min = min(h.Min, h.Buckets[i].Count)
				h.Max = max(h.Max, h.Buckets[i].Count)
				break
			}
		}
	}

	return h
}

// Create creates an histogram partioning input over `bins` buckets
func Create[T Indexable](bins int, input []T) Histogram[T] {
	if len(input) == 0 || bins == 0 {
		return Histogram[T]{}
	}

	sort.SliceStable(input, func(i, j int) bool {
		return input[i] < input[j]
	})

	minValue := input[0]
	maxValue := input[len(input)-1]

	scale := float64(maxValue-minValue) / float64(bins)

	if minValue == maxValue {
		bins = 1
	}

	buckets := make([]Bucket[T], bins)
	for i := 0; i < bins; i++ {
		buckets[i] = Bucket[T]{
			Min: T(float64(i)*scale) + minValue,
			Max: T(float64(i+1)*scale) + minValue,
		}
	}

	return newHistogram(buckets, input)
}

// CreateRanged creates a histogram by custom range
func CreateRanged[T Indexable](minValue, maxValue, interval T, input []T) Histogram[T] {
	if interval == 0 {
		return Histogram[T]{}
	}

	sort.SliceStable(input, func(i, j int) bool {
		return input[i] < input[j]
	})

	bins := 1
	if maxValue != minValue {
		bins = int(float64(maxValue-minValue) / float64(interval))
	}

	buckets := make([]Bucket[T], bins)

	for i := 0; i < bins; i++ {
		buckets[i] = Bucket[T]{
			Min: minValue + (interval * T(i)),
			Max: min(minValue+(interval*T(i+1)), maxValue),
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
		return input[i] < input[j]
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

func min[T Indexable](a T, b T) T {
	if a < b {
		return a
	}
	return b
}

func max[T Indexable](a T, b T) T {
	if a < b {
		return b
	}
	return a
}
