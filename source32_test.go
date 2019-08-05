package pcg

import "testing"

var sums32 = []struct {
	state    uint64 // PCG seed value for state
	sequence uint64 // PCG seed value for sequence
	count    int    // number of values to sum
	sum      uint64 // sum of the first count values
}{
	{1, 1, 1, 3380776849},
	{1, 1, 10, 23876797252},
	{1, 1, 100, 214946922057},
	{1, 1, 1000, 2103475364494},
	{1, 1, 10000, 21493472477812},
}

func Test32(t *testing.T) {
	for i, a := range sums32 {
		src := NewSource32(a.state, a.sequence)

		sum := uint64(0)
		for j := 0; j < a.count; j++ {
			sum += uint64(src.Uint32())
		}

		if sum != a.sum {
			t.Errorf("#%d, sum of first %d values = %d; want %d", i, a.count, sum, a.sum)
		}
	}
}
