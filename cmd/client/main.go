package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/et-nik/word-of-wisdom/internal/domain"
	"github.com/et-nik/word-of-wisdom/internal/solver"
	"github.com/pkg/errors"
)

func main() {
	tm := time.Now()

	host := flag.String("host", "", "server host")

	port := flag.String("port", "9100", "server port")
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}

	err = conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Fatalf("failed to set read deadline: %v", err)
	}

	reader := bufio.NewReader(conn)

	requestID, ch, err := processChallenge(conn, reader)
	if err != nil {
		log.Fatalf("failed to process challenge: %v", err)
	}

	quote, err := processWordOfWisdom(requestID, ch, conn, reader)
	if err != nil {
		log.Fatalf("failed to process word of wisdom: %v", err)
	}

	log.Println("Time elapsed:", time.Since(tm))

	// This is cli app, so we can use fmt.Println and fmt.Printf.
	fmt.Println("\n\n-----------------------------") //nolint:forbidigo
	fmt.Print(quote)                                 //nolint:forbidigo
	fmt.Println("-----------------------------")     //nolint:forbidigo
}

func processChallenge(w io.Writer, r *bufio.Reader) (string, domain.Challenge, error) {
	_, err := w.Write([]byte("challenge\n"))
	if err != nil {
		return "", domain.Challenge{}, errors.Wrap(err, "failed to write to server")
	}

	b, err := r.ReadBytes('\n')
	if err != nil {
		return "", domain.Challenge{}, errors.Wrap(err, "failed to read from server")
	}

	requestID, ch, err := parseChallenge(b)
	if err != nil {
		return "", domain.Challenge{}, errors.Wrap(err, "failed to parse challenge")
	}

	return requestID, ch, nil
}

func processWordOfWisdom(requestID string, ch domain.Challenge, w io.Writer, r *bufio.Reader) (string, error) {
	s := solver.New()
	sol, err := s.Solve(ch)
	if err != nil {
		return "", errors.Wrap(err, "failed to solve challenge")
	}

	_, err = w.Write(
		[]byte(fmt.Sprintf(
			"words-of-wisdom %s %d %x\n", requestID, sol.Nonce, sol.Hash,
		)),
	)
	if err != nil {
		return "", errors.Wrap(err, "failed to write to server")
	}

	b, err := r.ReadBytes('\n')
	if err != nil {
		return "", errors.Wrap(err, "failed to read from server")
	}

	if len(b) <= 2 {
		return "", errors.New("server responded with empty message")
	}

	if b[0] == '0' {
		return "", fmt.Errorf("server responded with error: %s", b[2:]) //nolint:goerr113
	}

	return string(b), nil
}

func parseChallenge(b []byte) (string, domain.Challenge, error) {
	var n, k int
	var seedStr string

	var requestID string

	_, err := fmt.Sscanf(string(b), "%s %d %d %s", &requestID, &n, &k, &seedStr)
	if err != nil {
		return "", domain.Challenge{}, errors.Wrap(err, "failed to parse challenge")
	}

	seed, err := hex.DecodeString(seedStr)
	if err != nil {
		return "", domain.Challenge{}, errors.Wrap(err, "failed to decode seed")
	}

	return requestID, domain.Challenge{
		N:    n,
		K:    k,
		Seed: seed,
	}, nil
}
