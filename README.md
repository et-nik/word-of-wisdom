# Word of Wisdom

[![Coverage Status](https://coveralls.io/repos/github/et-nik/word-of-wisdom/badge.svg?branch=main)](https://coveralls.io/github/et-nik/word-of-wisdom?branch=main)

Protected TCP Server from DOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol is used.

Used Memory bound PoW algorithm [Equihash](https://www.cryptolux.org/index.php/Equihash).

## Usage

### Build

```shell
make build-all
```

### Run

#### Server

Command run the server in Docker container with default parameters

```shell
make run-server
```

#### Client

Command run the client in Docker container with default parameters

```shell
make run-client
```

### Clean

Clean up all containers and images

```shell
make clean
```

## Selection of an algorithm

There are a lot of different algorithms for PoW, each with its own pros and cons.

I've examined a few:
* **Hashcash**. Easy to use, nice documentation. I decided to not use it due to is very sensitive to hardware.
* **Guided tour puzzle**. I liked the idea of this algorighm, but this is hard to maintain. THis algorithm is required extra services.
* **Equihash**. I considered this algorithm, it's memory-bound.

I decided to use Equihash-based algorithm bacause it is memory hard Proof-of-Work with fast verification.
Memory-bound Proof of Work (PoW) is selected due to its lower sensitivity to hardware variations, making it suitable for both low and high-end hardware. Additionally, the algorithm's performance is anticipated to be less affected by hardware advancements.

Difficulty of this PoW can be configured with two parameters _k_ and _n_. 
_n_ is the width (number of bits) of the generalized birthday problem.
_k_ is the length of the generalized birthday problem, should be small integer.


You can configurate this parameters in environment variables:
* **DIFFICULTY_WIDTH** represents _n_ 
* **DIFFICULTY_LENGTH** represents _k_

#### Recommended k and n values

| N  | K | Solve time |
|----|---|------------|
| 30 | 2 | 2 ms       |
| 40 | 3 | 30 ms      |
| 52 | 3 | 50 ms      |
| 56 | 3 | 60 ms      |
| 60 | 3 | 90 ms      |
| 70 | 3 | 400 ms     |
| 80 | 3 | 3.8 sec    |


### Performance & Benchmarks

#### Verification

Verification perform very fast, it's less than 100ns.

```
goos: darwin
goarch: arm64
pkg: github.com/et-nik/word-of-wisdom/internal/challenger
BenchmarkVerifier_Verify
BenchmarkVerifier_Verify-10    	12917673	        92.06 ns/op	     112 B/op	       4 allocs/op
```

#### Solver

```

goos: darwin
goarch: arm64
pkg: github.com/et-nik/word-of-wisdom/internal/solver
Benchmark
Benchmark/fast_challenge_10_2
Benchmark/fast_challenge_10_2-10         	  119826	      8824 ns/op	    7672 B/op	     187 allocs/op
Benchmark/challenge_50_2
Benchmark/challenge_50_2-10              	      13	  84333872 ns/op	55575446 B/op	 1245211 allocs/op
Benchmark/challenge_70_3
Benchmark/challenge_70_3-10              	       5	 206689875 ns/op	149424361 B/op	 3276875 allocs/op
```