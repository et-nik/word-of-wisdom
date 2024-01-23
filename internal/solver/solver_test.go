package solver_test

import (
	"testing"

	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/et-nik/word-of-wisdom/internal/solver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSolver(t *testing.T) {
	tests := []struct {
		name      string
		challenge domain.Challenge
		want      domain.Solution
	}{
		{
			name: "simple challenge",
			challenge: domain.Challenge{
				N:    4,
				K:    2,
				Seed: []byte{0x01, 0x02, 0x03, 0x04},
			},
			want: domain.Solution{
				Nonce: 15,
				Hash: []byte{
					0x0, 0x0, 0x0, 0x0,
					0x1, 0x0, 0x0, 0x0,
					0x2, 0x0, 0x0, 0x0,
					0x3, 0x0, 0x0, 0x0,
				},
			},
		},
		{
			name: "other challenge",
			challenge: domain.Challenge{
				N:    60,
				K:    3,
				Seed: []byte{50, 34, 234, 1, 3, 33, 192},
			},
			want: domain.Solution{
				Nonce: 2,
				Hash: []uint8{
					0xd1, 0x5d, 0x0, 0x0,
					0x5f, 0x69, 0x0, 0x0,
					0x76, 0x69, 0x0, 0x0,
					0x51, 0x7e, 0x0, 0x0,
					0x7c, 0x95, 0x0, 0x0,
					0xc0, 0xbb, 0x0, 0x0,
					0xe4, 0xc4, 0x0, 0x0,
					0x9c, 0xcd, 0x0, 0x0,
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// ARRANGE
			s := solver.New()

			// ACT
			result, err := s.Solve(test.challenge)

			// ASSERT
			require.NoError(t, err)
			assert.Equal(t, test.want, result)
		})
	}
}

func Benchmark(b *testing.B) {
	benchmarks := []struct {
		name      string
		challenge domain.Challenge
	}{
		{
			name: "fast challenge 10 2",
			challenge: domain.Challenge{
				N:    10,
				K:    2,
				Seed: []byte{0x01, 0x02, 0x03, 0x04},
			},
		},
		{
			name: "challenge 50 2",
			challenge: domain.Challenge{
				N:    50,
				K:    2,
				Seed: []byte{0x01, 0x02, 0x03, 0x04},
			},
		},
		{
			name: "challenge 70 3",
			challenge: domain.Challenge{
				N:    70,
				K:    3,
				Seed: []byte{0x01, 0x02, 0x03, 0x04},
			},
		},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			s := solver.New()

			for i := 0; i < b.N; i++ {
				_, _ = s.Solve(bm.challenge)
			}
		})
	}
}
