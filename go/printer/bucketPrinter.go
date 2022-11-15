package printer

import (
	"io"
	"math"
	"strings"
	"text/tabwriter"
	"text/template"

	histogram "github.com/huffduff/histogram/go"
)

type bucketPrinter[T histogram.Indexable] struct {
	histogram.Bucket[T]
	h histogram.Histogram[T]
}

func (p bucketPrinter[T]) Total() int {
	return p.h.Count
}

func (p bucketPrinter[T]) Pct() float64 {
	total := p.Total()
	if total == 0 {
		return 0
	}
	return float64(p.Count) / float64(total) * 100.0
}

func (p bucketPrinter[T]) scale() float64 {
	if p.Total() == 0 {
		return 0
	}
	if p.h.Min == p.h.Max {
		return 1
	}
	return float64(p.Count-p.h.Min) / float64(p.h.Max-p.h.Min)
}

func (p bucketPrinter[T]) Bar(width float64) string {
	size := p.scale() * width
	decimalf := (size - math.Floor(size)) * 10.0
	decimali := math.Floor(decimalf)
	charIdx := int(decimali / 10.0 * 8.0)
	return strings.Repeat("█", int(size)) + string(rune(9615-charIdx))
}

// Fprint writes a histogram to an io.Writer using the format strings provided
// The format strings use the text/template syntax and write to a text/tabwriter
// multiple format strings are joined in order using \t to create aligned columns
// and each line automatically gets a \n, so no need to add an EOL.
// Values passed to the template:
//
//	.Total : the total number of records in the histogram (int)
//	.Pct   : percentage represesented by a bucket value relative to .Total (float64)
//	.Min   : lower bound of the bucket, inclusive
//	.Max   : higher bound of the bucket, exclusive except on the final bucket
//	.Count : count of values in the bucket
//
// [formatting.go](.formatting.go) contains some helper functions for common formats.
// If no formats are passed, it uses `IntFormat(5)` by default.
func Fprint[T histogram.Indexable](w io.Writer, h histogram.Histogram[T], formats ...string) error {
	if len(formats) == 0 {
		formats = IntFormat(5)
	}
	return fprintf(w, h, formats)
}

func fprintf[T histogram.Indexable](w io.Writer, h histogram.Histogram[T], formats []string) error {
	tabw := tabwriter.NewWriter(w, 2, 2, 2, byte(' '), 0)

	line, err := template.New("line").Parse(strings.Join(formats, "\t") + "\n")
	if err != nil {
		return err
	}

	for _, v := range h.Buckets {
		b := bucketPrinter[T]{v, h}
		err := line.Execute(tabw, b)
		if err != nil {
			return err
		}
	}

	return tabw.Flush()
}
