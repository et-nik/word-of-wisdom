package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	config2 "github.com/et-nik/word-of-wisdom/internal/config"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type handler func(ctx context.Context, msg Message, writer io.Writer)

type Server struct {
	config   *config2.Config
	handlers map[string]handler
	listener net.Listener
}

func NewServer(config *config2.Config) *Server {
	return &Server{
		config:   config,
		handlers: make(map[string]handler),
	}
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server is running")

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf("%s:%d", s.config.ListenAddr, s.config.ListenPort))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	log.Printf("Server is listening on %s:%d\n", s.config.ListenAddr, s.config.ListenPort)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return nil
			}
			log.Println("failed to accept connection, err:", err)
			continue
		}

		go s.handleConnection(ctx, conn)
	}
}

// RegisterHandler registers a handler for a given message type
// If a handler for a given message type is already registered, it will be overwritten.
func (s *Server) RegisterHandler(messageType string, h handler) {
	s.handlers[messageType] = h
}

func (s *Server) Stop() error {
	var merr error

	err := s.listener.Close()
	merr = multierr.Append(merr, err)

	return merr
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("failed to close connection, err:", err)
		}
	}(conn)

	err := conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		log.Println("failed to set read deadline, err:", err)
		return
	}

	reader := bufio.NewReader(conn)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("failed to read bytes, err:", err)
			}
			return
		}
		log.Printf("request: %s", bytes)

		msg, err := parseMessage(bytes)
		if err != nil {
			log.Println("failed to parse message, err:", err)
			return
		}

		h, ok := s.handlers[msg.Type]
		if !ok {
			log.Printf("handler for message type %s not found", msg.Type)
			return
		}

		h(ctx, msg, conn)
		_, err = conn.Write([]byte("\n"))
		if err != nil {
			log.Println("failed to write response, err:", err)
			return
		}
	}
}
