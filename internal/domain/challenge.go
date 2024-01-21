package domain

type Challenge struct {
	Algorithm string

	// The width (number of bits) of the generalized birthday problem, integer divisible by (K+1)
	N int

	// The length of the generalized birthday problem, small integer
	K int

	Seed []byte
}

type Solution struct {
	Nonce int
	Hash  []byte
}
