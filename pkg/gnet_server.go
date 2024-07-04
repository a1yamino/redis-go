package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
)

type GServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	net       string
	addr      string
	multicore bool
	handlers  map[string]CommandHandler
	pool      *goroutine.Pool
}

func NewGServer(net, addr string, multicore bool, p *goroutine.Pool) *GServer {
	handlers := make(map[string]CommandHandler)

	for cmd, handler := range defaultHandlers {
		handlers[cmd] = handler
	}

	return &GServer{
		net:       net,
		addr:      addr,
		multicore: multicore,
		handlers:  handlers,
		pool:      p,
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
	c := NewGConn(conn)

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

	_ = s.pool.Submit(func() {
		handler.call(c.Writer, req[1:])
	})

	return gnet.None
}

type gConn struct {
	gnet.Conn
	Writer     *gWriter
	Reader     *Reader
	remoteAddr string
}

func NewGConn(conn gnet.Conn) *gConn {
	return &gConn{
		Conn:       conn,
		Writer:     NewGWriter(conn),
		Reader:     NewReader(conn),
		remoteAddr: conn.RemoteAddr().String(),
	}
}

type gWriter struct {
	gnet.Writer
}

var nilCallback gnet.AsyncCallback

func NewGWriter(w gnet.Writer) *gWriter {
	return &gWriter{w}
}

func (w *gWriter) WriteSimpleString(str string) error {
	err := w.AsyncWrite([]byte("+"+str+"\r\n"), nilCallback)
	return err
}

func (w *gWriter) WriteError(str string) error {
	err := w.AsyncWrite([]byte("-"+str+"\r\n"), nilCallback)
	return err
}

func (w *gWriter) WriteInteger(i int) error {
	err := w.AsyncWrite([]byte(":"+strconv.Itoa(i)+"\r\n"), nilCallback)
	return err
}

func (w *gWriter) WriteBulkString(str string) error {
	err := w.AsyncWrite([]byte("$"+strconv.Itoa(len(str))+"\r\n"+str+"\r\n"), nilCallback)
	return err
}

func (w *gWriter) WriteNull() error {
	err := w.AsyncWrite([]byte("$-1\r\n"), nilCallback)
	return err
}

func (w *gWriter) WriteArray(value Value) error {
	err := w.AsyncWrite([]byte("*"+strconv.Itoa(len(value.array))+"\r\n"), nilCallback)

	if err != nil {
		return err
	}

	for _, v := range value.array {
		switch v.typ {
		case STRING:
			err = w.WriteSimpleString(v.str)
		case ERROR:
			err = w.WriteError(v.str)
		case INTEGER:
			err = w.WriteInteger(v.integer)
		case BULK:
			err = w.WriteBulkString(v.str)
		}
	}
	return err
}
