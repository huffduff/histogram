package printer

import "fmt"

// Pct generates a .Pct value using the supplied Printf format
func Pct(format string) string {
	return `{{printf "` + format + `%%" .Pct }}`
}

// Count generates a .Count value using the supplied Printf format
func Count(format string) string {
	return `{{printf "` + format + `" .Count}}`
}

// Pct generates a .Pct value using the `pos` Printf format if
// the value is not zero, `pos` on a zero value.
// useful for hiding zero values
func CountIf(pos string, neg string) string {
	return `{{if (eq .Count 0)}}` + neg + `{{else}}{{printf "` + pos + `" .Count}}{{ end }}`
}

// Range generates a .Min and .Max values using the supplied Printf format
func Range(format string) string {
	return `{{printf "` + format + `" .Min .Max}}`
}

// Bar generates the actual .Bar with the max length provided
func Bar(length int) string {
	return fmt.Sprintf("{{.Bar %d}}", length)
}

// FloatFormat generates a tabulated row with the bar length provided
// when the value is a float
func FloatFormat(length int) []string {
	return []string{
		Range("%.3g-%.3g"),
		Pct("%.4g"),
		Bar(length),
		CountIf("% 4d", "   -"),
	}
}

// IntFormat generates a tabulated row with the bar length provided
// when the value is an int
func IntFormat(length int) []string {
	return []string{
		Range("%d-%d"),
		Pct("%.4g"),
		Bar(length),
		CountIf("% 4d", "   -"),
	}
}

// StringFormat generates a tabulated row with the bar length provided
// when the value implements fmt.Stringer (time.Duration for example)
func StringFormat(length int) []string {
	return []string{
		Range("%s"),
		Pct("%.4g"),
		Bar(length),
		CountIf("% 4d", "   -"),
	}
}
