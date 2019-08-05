package pcg

const mask = ^uint64(0) >> 1

type Source64 struct {
	lo *Source32
	hi *Source32
}

func NewSource64(state1, state2, sequence1, sequence2 uint64) *Source64 {
	if sequence1&mask == sequence2&mask {
		sequence2 = ^sequence2
	}

	return &Source64{
		NewSource32(state1, sequence1),
		NewSource32(state2, sequence2),
	}
}

func (src *Source64) Uint64() uint64 {
	return uint64(src.hi.Uint32())<<32 | uint64(src.lo.Uint32())
}
