package printer

import "fmt"

func Pct(format string) string {
	return `{{printf "` + format + `%%" .Pct }}`
}

func Count(format string) string {
	return `{{printf "` + format + `" .Count}}`
}

func CountIf(pos string, neg string) string {
	return `{{if (eq .Count 0)}}` + neg + `{{else}}{{printf "` + pos + `" .Count}}{{ end }}`
}

func Range(format string) string {
	return `{{printf "` + format + `" .Min .Max}}`
}

func Bar(length int) string {
	return fmt.Sprintf("{{.Bar %d}}", length)
}

func StandardFormat(length int) []string {
	return []string{
		Range("%.3g-%.3g"),
		Pct("%.4g"),
		Bar(length),
		CountIf("% 4d", "   -"),
		// Count("% 4d"),
	}
}

func StringFormat(length int) []string {
	return []string{
		Range("%s"),
		Pct("%.4g"),
		Bar(length),
		CountIf("% 4d", "   -"),
	}
}
