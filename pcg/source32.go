package pcg

const multiplier = 6364136223846793005

type Source32 struct {
	state     uint64
	increment uint64
}

func NewSource32(state, sequence uint64) *Source32 {
	increment := (sequence << 1) | 1
	state = (state+increment)*multiplier + increment

	return &Source32{state, increment}
}

func (src *Source32) Uint32() uint32 {
	xor := uint32(((src.state >> 18) ^ src.state) >> 27)
	rot := uint32(src.state >> 59)

	src.state = src.state*multiplier + src.increment

	return (xor >> rot) | (xor << ((-rot) & 31))
}
