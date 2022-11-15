package histogram

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	t.Run("floats", func(t *testing.T) {
		data := []float64{
			0.1,
			0.2, 0.21, 0.22, 0.22,
			0.3,
			0.4,
			0.5, 0.51, 0.52, 0.53, 0.54, 0.55, 0.56, 0.57, 0.58,
			0.6,
			// 0.7 is empty
			0.8,
			0.9,
			1.0,
		}
		hist := Create(9, data)

		if len(hist.Buckets) != 9 {
			t.Errorf("%d != %d", len(hist.Buckets), 9)
		}

		idx, err := hist.Index(1.0)
		if err != nil {
			t.Error(err)
		}
		if idx != 8 {
			t.Errorf("bucket %d != %d", idx, 8)
		}
	})

	t.Run("duratation", func(t *testing.T) {
		data := []time.Duration{
			time.Millisecond * 100,
			time.Millisecond * 200,
			time.Millisecond * 210,
			time.Millisecond * 220,
			time.Millisecond * 221,
			time.Millisecond * 222,
			time.Millisecond * 223,
			time.Millisecond * 300,
			time.Millisecond * 400,
			time.Millisecond * 500,
			time.Millisecond * 510,
			time.Millisecond * 520,
			time.Millisecond * 530,
			time.Millisecond * 540,
			time.Millisecond * 550,
			time.Millisecond * 560,
			time.Millisecond * 570,
			time.Millisecond * 580,
			time.Millisecond * 600,
			// no 0.7s
			time.Millisecond * 800,
			time.Millisecond * 900,
			time.Millisecond * 1000,
		}

		hist := Create(9, data)

		if len(hist.Buckets) != 9 {
			t.Errorf("%d != %d", len(hist.Buckets), 9)
		}

		idx, err := hist.Index(time.Millisecond * 1000)
		if err != nil {
			t.Error(err)
		}
		if idx != 8 {
			t.Errorf("bucket %d != %d", idx, 8)
		}
	})
}

func TestCreateRanged(t *testing.T) {
	data := []float64{
		1,
		2, 2.1, 2.5, 2.8,
		3,
		5,
		5.5,
		6.5,
		6.6,
		7,
		10,
	}

	t.Run("includeRange", func(t *testing.T) {
		hist := CreateRanged(0.5, 20, 2, data)

		if hist.Count != len(data) {
			t.Errorf("histogram count: %d != %d", hist.Count, len(data))
		}
		if hist.Buckets[3].Count != 3 {
			t.Errorf("histogram buckets: %d != %d", hist.Buckets[3].Count, 3)
		}
	})
	t.Run("crossRange", func(t *testing.T) {
		hist := CreateRanged(6, 20, 2, data)

		if hist.Count != 4 {
			t.Errorf("histogram count: %d != %d", hist.Count, 4)
		}
	})
}

func TestByZeroRange(t *testing.T) {

	data := []float64{
		-1,
		1,
	}

	hist := CreateRanged(6, 20, 2, data)

	if hist.Count != 0 {
		t.Errorf("histogram count: %d != %d", hist.Count, 0)
	}
	if len(hist.Buckets) != 7 {
		t.Errorf("buckets count: %d != %d", len(hist.Buckets), 7)
	}
}

func TestCreateLog(t *testing.T) {
	data := []float64{
		1,
		2, 2.1, 2.5, 2.8,
		3,
		5,
		5.5,
		6.5,
		6.6,
		7,
		10,
		100,
	}
	t.Run("binary", func(t *testing.T) {
		hist := CreateLog(2, data)
		idx, err := hist.Index(10)
		if err != nil {
			t.Error(err)
		}
		if idx != 3 {
			t.Errorf("%d != %d", idx, 3)
		}
	})
	t.Run("decimal", func(t *testing.T) {
		hist := CreateLog(10, data)
		idx, err := hist.Index(10)
		if err != nil {
			t.Error(err)
		}
		if idx != 1 {
			t.Errorf("%d != %d", idx, 1)
		}
	})
}
