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
				Nonce: 2,
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
				Nonce: 13,
				Hash: []uint8{
					0xde, 0x0, 0x0, 0x0,
					0x73, 0x26, 0x0, 0x0,
					0x3c, 0x42, 0x0, 0x0,
					0xdc, 0x45, 0x0, 0x0,
					0xc7, 0x87, 0x0, 0x0,
					0x4d, 0x93, 0x0, 0x0,
					0xc9, 0xe2, 0x0, 0x0,
					0x18, 0xfe, 0x0, 0x0,
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
			name: "fast challenge",
			challenge: domain.Challenge{
				N:    10,
				K:    2,
				Seed: []byte{0x01, 0x02, 0x03, 0x04},
			},
		},
		{
			name: "challenge",
			challenge: domain.Challenge{
				N:    50,
				K:    2,
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
