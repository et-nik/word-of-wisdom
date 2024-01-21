package equihash

import (
	"encoding/binary"
	"sort"

	"github.com/dchest/blake2b"
)

const (
	seedLength = 16
	maxN       = 200
	listLength = 4
	maxNonce   = 1 << 30

	forkMultiplier = 3
)

type Input uint32

type Proof struct {
	n      uint32
	k      uint32
	seed   [seedLength]uint32
	nonce  uint32
	inputs []Input
}

func NewProof(
	n int,
	k int,
	seed []byte,
	nonce int,
	inputs []byte,
) *Proof {
	var s [16]uint32
	for i, c := range BytesToUint32Array(seed, binary.LittleEndian) {
		if i > 15 {
			break
		}

		s[i] = c
	}

	inp := make([]Input, len(inputs))
	for i, c := range BytesToUint32Array(inputs, binary.LittleEndian) {
		inp[i] = Input(c)
	}

	return &Proof{
		n:      uint32(n),
		k:      uint32(k),
		seed:   s,
		nonce:  uint32(nonce),
		inputs: inp,
	}
}

func (p *Proof) Inputs() []Input {
	return p.inputs
}

func (p *Proof) InputsBytes() []byte {
	inputs := make([]uint32, len(p.inputs))
	for i := range p.inputs {
		inputs[i] = uint32(p.inputs[i])
	}

	return Uint32ArrayToBytes(inputs, binary.LittleEndian)
}

func (p *Proof) Nonce() uint32 {
	return p.nonce
}

type fork struct {
	ref1 uint32
	ref2 uint32
}

type tuple struct {
	blocks    []uint32
	reference uint32
}

type Equihash struct {
	n          uint32
	k          uint32
	seed       [seedLength]uint32
	nonce      uint32
	tupleList  [][]tuple
	filledList []uint32
	solutions  []Proof
	forks      [][]fork
}

func New(n int, k int, seed []byte) *Equihash {
	var s [16]uint32
	for i, c := range BytesToUint32Array(seed, binary.LittleEndian) {
		if i > 15 {
			break
		}

		s[i] = c
	}

	return &Equihash{
		n:     uint32(n),
		k:     uint32(k),
		seed:  s,
		nonce: 1,
	}
}

func (eq *Equihash) initializeMemory() {
	tupleN := uint32(1) << (eq.n / (eq.k + 1))
	defaultTuple := tuple{blocks: make([]uint32, eq.k)}
	defTuples := make([]tuple, listLength)
	for i := range defTuples {
		defTuples[i] = defaultTuple
	}
	eq.tupleList = make([][]tuple, tupleN)
	for i := range eq.tupleList {
		eq.tupleList[i] = append([]tuple(nil), defTuples...)
	}
	eq.filledList = make([]uint32, tupleN)
	eq.solutions = make([]Proof, 0)
	eq.forks = make([][]fork, 0)
}

func (eq *Equihash) fillMemory(length uint32) {
	input := make([]uint32, seedLength+2)
	copy(input[:seedLength], eq.seed[:])
	input[seedLength] = eq.nonce
	buf := make([]uint32, maxN/4)
	h := blake2b.New256()
	for i := uint32(0); i < length; i++ {
		input[seedLength+1] = i

		h.Reset()
		h.Write(Uint32ArrayToBytes(input[:], binary.LittleEndian))
		buf = BytesToUint32Array(h.Sum(nil), binary.LittleEndian)

		index := buf[0] >> (32 - eq.n/(eq.k+1))
		count := eq.filledList[index]
		if count < listLength {
			for j := uint32(1); j < (eq.k + 1); j++ {
				eq.tupleList[index][count].blocks[j-1] = buf[j] >> (32 - eq.n/(eq.k+1))
			}
			eq.tupleList[index][count].reference = i
			eq.filledList[index]++
		}
	}
}

func (eq *Equihash) resolveCollisions(store bool) {
	tableLength := len(eq.tupleList)
	maxNewCollisions := len(eq.tupleList) * forkMultiplier
	newBlocks := len(eq.tupleList[0][0].blocks) - 1
	newForks := make([]fork, maxNewCollisions)
	tableRow := make([]tuple, listLength)
	for i := range tableRow {
		tableRow[i] = tuple{blocks: make([]uint32, newBlocks)}
	}
	collisionList := make([][]tuple, tableLength)
	for i := range collisionList {
		collisionList[i] = append([]tuple(nil), tableRow...)
	}
	newFilledList := make([]uint32, tableLength)
	newColls := uint32(0)
	for i := uint32(0); i < uint32(tableLength); i++ {
		for j := uint32(0); j < eq.filledList[i]; j++ {
			for m := j + 1; m < eq.filledList[i]; m++ {
				newIndex := eq.tupleList[i][j].blocks[0] ^ eq.tupleList[i][m].blocks[0]
				newFork := fork{ref1: eq.tupleList[i][j].reference, ref2: eq.tupleList[i][m].reference}
				if store {
					if newIndex == 0 {
						solutionInputs := eq.resolveTree(newFork)
						eq.solutions = append(
							eq.solutions,
							Proof{n: eq.n, k: eq.k, seed: eq.seed, nonce: eq.nonce, inputs: solutionInputs},
						)
					}
				} else {
					if newFilledList[newIndex] < listLength && newColls < uint32(maxNewCollisions) {
						for l := uint32(0); l < uint32(newBlocks); l++ {
							collisionList[newIndex][newFilledList[newIndex]].blocks[l] = eq.tupleList[i][j].blocks[l+1] ^ eq.tupleList[i][m].blocks[l+1]
						}
						newForks[newColls] = newFork
						collisionList[newIndex][newFilledList[newIndex]].reference = newColls
						newFilledList[newIndex]++
						newColls++
					}
				}
			}
		}
	}
	eq.forks = append(eq.forks, newForks)
	eq.tupleList, collisionList = collisionList, eq.tupleList
	eq.filledList, newFilledList = newFilledList, eq.filledList
}

func (eq *Equihash) resolveTreeByLevel(fork fork, level uint32) []Input {
	if level == 0 {
		return []Input{Input(fork.ref1), Input(fork.ref2)}
	}
	v1 := eq.resolveTreeByLevel(eq.forks[level-1][fork.ref1], level-1)
	v2 := eq.resolveTreeByLevel(eq.forks[level-1][fork.ref2], level-1)
	return append(v1, v2...)
}

func (eq *Equihash) resolveTree(fork fork) []Input {
	return eq.resolveTreeByLevel(fork, uint32(len(eq.forks)))
}

func (eq *Equihash) FindProof() Proof {
	eq.nonce = 1
	for eq.nonce < maxNonce {
		eq.nonce++
		eq.initializeMemory()
		eq.fillMemory(4 << (eq.n/(eq.k+1) - 1))
		for i := uint32(1); i <= eq.k; i++ {
			toStore := (i == eq.k)
			eq.resolveCollisions(toStore)
		}
		for i := range eq.solutions {
			vec := eq.solutions[i].inputs
			sort.Slice(vec, func(i, j int) bool { return vec[i] < vec[j] })
			dup := false
			for k := range vec[:len(vec)-1] {
				if vec[k] == vec[k+1] {
					dup = true
				}
			}
			if !dup {
				return eq.solutions[i]
			}
		}
	}
	return Proof{n: eq.n, k: eq.k, seed: eq.seed, nonce: eq.nonce, inputs: []Input{}}
}

func (p *Proof) Test() bool {
	input := make([]uint32, seedLength+2)
	copy(input[:seedLength], p.seed[:])
	input[seedLength] = p.nonce
	buf := make([]uint32, maxN/4)
	blocks := make([]uint32, p.k+1)

	h := blake2b.New256()

	for i := range p.inputs {
		input[seedLength+1] = uint32(p.inputs[i])

		h.Reset()
		h.Write(Uint32ArrayToBytes(input[:], binary.LittleEndian))
		buf = BytesToUint32Array(h.Sum(nil), binary.LittleEndian)

		for j := uint32(0); j < (p.k + 1); j++ {
			//select j-th block of n/(k+1) bits
			blocks[j] ^= buf[0] >> (32 - p.n/(p.k+1))
		}
	}

	for _, block := range blocks {
		if block != 0 {
			return false
		}
	}

	return true
}
