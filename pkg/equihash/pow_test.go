package equihash

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquihash(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		k        int
		seed     []uint32
		expected Proof
	}{
		{
			name: "52_3_1",
			n:    52,
			k:    3,
			seed: []uint32{1, 1, 1, 1},
			expected: *NewProof(
				52,
				3,
				Uint32ArrayToBytes([]uint32{1, 1, 1, 1}, binary.LittleEndian),
				2,
				Uint32ArrayToBytes([]uint32{0x611, 0x1626, 0x1c37, 0x20cb, 0x241d, 0x30d7, 0x3811, 0x395c}, binary.LittleEndian),
			),
		},
		{
			name: "30_2_5",
			n:    52,
			k:    3,
			seed: []uint32{5, 5, 5, 5, 5, 5, 5, 5},
			expected: *NewProof(
				52,
				3,
				Uint32ArrayToBytes([]uint32{5, 5, 5, 5, 5, 5, 5, 5}, binary.LittleEndian),
				3,
				Uint32ArrayToBytes([]uint32{0x1238, 0x12e6, 0x1500, 0x2aae, 0x2db4, 0x3852, 0x393f, 0x3d92}, binary.LittleEndian),
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// ARRANGE
			eq := New(test.n, test.k, Uint32ArrayToBytes(test.seed, binary.LittleEndian))

			// ACT
			proof := eq.FindProof()

			// ASSERT
			assert.Equal(t, test.expected, proof)
		})
	}
}

func TestProof_Test(t *testing.T) {
	tests := []struct {
		name     string
		proof    Proof
		expected bool
	}{
		{
			name: "52_3_1",
			proof: *NewProof(
				52,
				3,
				Uint32ArrayToBytes([]uint32{1, 1, 1, 1}, binary.LittleEndian),
				2,
				Uint32ArrayToBytes([]uint32{0x611, 0x1626, 0x1c37, 0x20cb, 0x241d, 0x30d7, 0x3811, 0x395c}, binary.LittleEndian),
			),
			expected: true,
		},
		{
			name: "52_3_1_wrong",
			proof: *NewProof(
				52,
				3,
				Uint32ArrayToBytes([]uint32{1, 1, 1, 1}, binary.LittleEndian),
				2,
				Uint32ArrayToBytes([]uint32{0x610, 0x1626, 0x1c37, 0x20cb, 0x241d, 0x30d7, 0x3811, 0x395c}, binary.LittleEndian),
			),
			expected: false,
		},
		{
			name: "52_3_2_wrong",
			proof: *NewProof(
				52,
				3,
				Uint32ArrayToBytes([]uint32{2, 2, 2, 2}, binary.LittleEndian),
				2,
				Uint32ArrayToBytes([]uint32{0x611, 0x1626, 0x1c37, 0x20cb, 0x241d, 0x30d7, 0x3811, 0x395c}, binary.LittleEndian),
			),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.proof.Test()
			assert.Equal(t, test.expected, result)
		})
	}
}
