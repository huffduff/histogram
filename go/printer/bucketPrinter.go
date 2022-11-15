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

func Fprint[T histogram.Indexable](w io.Writer, h histogram.Histogram[T], formats ...string) error {
	if len(formats) == 0 {
		formats = StandardFormat(5)
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
