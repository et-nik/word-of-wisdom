package server

import (
	"bytes"
	"errors"
	"strings"
)

var errInvalidMessage = errors.New("invalid message")

type Message struct {
	Type    string
	Payload string
}

func parseMessage(msg []byte) (Message, error) {
	splitted := bytes.SplitN(msg, []byte(" "), 2)

	if len(splitted) == 0 {
		return Message{}, errInvalidMessage
	}

	if len(splitted) == 1 {
		return Message{
			Type: normalizeType(string(splitted[0])),
		}, nil
	}

	return Message{
		Type:    normalizeType(string(splitted[0])),
		Payload: strings.TrimSpace(string(splitted[1])),
	}, nil
}

func normalizeType(t string) string {
	return strings.ToLower(strings.TrimSpace(t))
}
