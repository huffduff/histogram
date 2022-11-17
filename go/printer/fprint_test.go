package printer

import (
	"os"

	histogram "github.com/huffduff/histogram/go"
)

func ExampleFprint() {
	data := []float64{
		0.1,
		0.2, 0.21, 0.22, 0.22,
		0.4,
		0.5, 0.51, 0.52, 0.53, 0.54, 0.55, 0.56, 0.57, 0.58, 0.59,
		0.6,
		// 0.7 is empty
		// 0.8,
		0.9,
		1.0,
		0.3, // intenionally out of order
		// 1000,
	}
	hist := histogram.Create(9, data)
	Fprint(os.Stdout, hist, FloatFormat(50)...)
	// Output:
	// 0.1-0.2  5%   █████                                                  1
	// 0.2-0.3  25%  █████████████████████████                              5
	// 0.3-0.4  0%   ▏                                                      -
	// 0.4-0.5  5%   █████                                                  1
	// 0.5-0.6  50%  ██████████████████████████████████████████████████    10
	// 0.6-0.7  5%   █████                                                  1
	// 0.7-0.8  0%   ▏                                                      -
	// 0.8-0.9  0%   ▏                                                      -
	// 0.9-1    10%  ██████████                                             2
}
