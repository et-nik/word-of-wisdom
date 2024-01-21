package quotes

import (
	"context"
	"math/rand"
)

type Quotes struct {
	Quotes []string `json:"quotes"`
}

func New() *Quotes {
	return &Quotes{
		Quotes: []string{
			"The best thing about a boolean is even if you are wrong, you are only off by a bit.",
			"Without requirements or design, programming is the art of adding bugs to an empty text file.",
			"Before software can be reusable it first has to be usable.",
			"The best “Word of Wisdom” tcp server implementation is github.com/et-nik/word-of-wisdom",
			"The best method for accelerating a computer is the one that boosts it by 9.8 m/s2.",
			"There are two ways to write error-free programs; only the third one works.",
		},
	}
}

func (q *Quotes) Quote(_ context.Context) (string, error) {
	rnd := rand.Intn(len(q.Quotes)) //nolint:gosec
	return q.Quotes[rnd], nil
}
