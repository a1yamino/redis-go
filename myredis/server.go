package myredis

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Conn struct {
	net.Conn
	Writer     *Writer
	Reader     *Reader
	remoteAddr string
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		Conn:       conn,
		Writer:     NewWriter(conn),
		Reader:     NewReader(conn),
		remoteAddr: conn.RemoteAddr().String(),
	}
}

type Config struct {
	EnableAof bool
	AofFile   string
}

type Server struct {
	sync.RWMutex
	handlers map[string]CommandHandler
	accpet   func(conn net.Conn) bool
	config   *Config
	Aof      *Aof
}

func NewServer(config *Config) *Server {
	handlers := make(map[string]CommandHandler)

	for cmd, handler := range defaultHandlers {
		handlers[cmd] = handler
	}

	s := &Server{
		handlers: handlers,
		config:   config,
	}

	if s.config != nil && s.config.EnableAof {
		bootstrapAof(s)
	}

	return s
}

func (s *Server) HandlerFunc(cmd string, handler CommandHandler) {
	s.Lock()
	defer s.Unlock()
	s.handlers[strings.ToUpper(cmd)] = handler
}

func (s *Server) AccpetFunc(f func(conn net.Conn) bool) {
	s.Lock()
	defer s.Unlock()
	s.accpet = f
}

func (s *Server) handleConn(conn net.Conn) error {
	c := NewConn(conn)
	defer conn.Close()
	s.Lock()
	accpet := s.accpet
	s.Unlock()
	if accpet != nil && !accpet(conn) {
		return nil
	}

	for {
		value, err := c.Reader.resp.Read()
		if err != nil {
			return err
		}

		if value.typ == 0 {
			continue
		}

		if !value.IsArray() {
			continue
		}

		req := value.Array()

		cmd := strings.ToUpper(req[0].String())

		s.RLock()
		handler, ok := s.handlers[cmd]
		s.RUnlock()

		if !ok {
			if err := c.Writer.WriteError("ERR unknown command '" + cmd + "'"); err != nil {
				return err
			}

			continue
		}

		if !handler.call(c, req[1:]) {
			return nil
		}

		if s.Aof != nil && handler.should_persist() {
			if err := s.Aof.Append(value); err != nil {
				return err
			}
		}
	}
}

func (s *Server) ListenAndServe(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go func() {
			err := s.handleConn(conn)
			defer conn.Close()

			if err != nil {
				if err == io.EOF {
					return
				}
				fmt.Println("handle connection failed: ", err)
			}
		}()
	}
}

func bootstrapAof(s *Server) {
	aof, err := NewAof(s.config.AofFile)
	if err != nil {
		panic(err)
	}

	s.Aof = aof

	aof.ReadValues(func(value Value) bool {
		cmds := value.Array()

		cmd := strings.ToUpper(cmds[0].String())

		s.RLock()
		handler, ok := s.handlers[cmd]
		s.RUnlock()

		if !ok {
			return true
		}

		// create fake connection with fake writer
		conn := &Conn{
			Writer: NewWriter(io.Discard),
		}

		handler.call(conn, cmds[1:])

		return true
	})
}
