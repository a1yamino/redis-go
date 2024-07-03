package pkg

import (
	"fmt"
	"strings"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
)

type GServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	net       string
	addr      string
	multicore bool
	handlers  map[string]CommandHandler
}

func NewGServer(net, addr string, multicore bool) *GServer {
	handlers := make(map[string]CommandHandler)

	for cmd, handler := range defaultHandlers {
		handlers[cmd] = handler
	}

	return &GServer{
		net:       net,
		addr:      addr,
		multicore: multicore,
		handlers:  handlers,
	}
}

func (s *GServer) ProtoAddr() string {
	return fmt.Sprintf("%s://%s", s.net, s.addr)
}

func (s *GServer) Multicore() bool {
	return s.multicore
}

func (s *GServer) OnBoot(eng gnet.Engine) gnet.Action {
	logging.Infof("running server on %s with multi-core=%t", fmt.Sprintf("%s://%s", s.net, s.addr), s.multicore)
	s.eng = eng
	return gnet.None
}

func (s *GServer) OnTraffic(conn gnet.Conn) gnet.Action {
	c := NewConn(conn)

	value, err := c.Reader.resp.Read()
	if err != nil {
		return gnet.Close
	}

	if value.typ == 0 {
		return gnet.None
	}

	if !value.IsArray() {
		return gnet.None
	}

	req := value.Array()

	cmd := strings.ToUpper(req[0].String())

	handler, ok := s.handlers[cmd]

	if !ok {
		if err := c.Writer.WriteError("ERR unknown command '" + cmd + "'"); err != nil {
			return gnet.Close
		}

		return gnet.None
	}

	if !handler.call(c, req[1:]) {
		return gnet.Close
	}

	return gnet.None
}
