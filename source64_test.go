package pcg

import "testing"

var sums64 = []struct {
	state1    uint64 // PCG seed value for state
	state2    uint64
	sequence1 uint64 // PCG seed value for sequence
	sequence2 uint64
	count     int    // number of values to sum
	sum       uint64 // sum of the first count values
}{
	{1, 1, 1, 2, 1, 1107300197865787281},
	{1, 1, 1, 2, 10, 8034187309725975364},
	{1, 1, 1, 2, 100, 14328956917741108809},
	{1, 1, 1, 2, 1000, 15814724732829753998},
	{1, 1, 1, 2, 10000, 8547922387302793844},
}

func Test64(t *testing.T) {
	for i, a := range sums64 {
		src := NewSource64(a.state1, a.state2, a.sequence1, a.sequence2)

		sum := uint64(0)
		for j := 0; j < a.count; j++ {
			sum += src.Uint64()
		}

		if sum != a.sum {
			t.Errorf("#%d, sum of first %d values = %d; want %d", i, a.count, sum, a.sum)
		}
	}
}
