package pkg

import (
	"testing"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
)

func TestGnetServer(t *testing.T) {
	p := goroutine.Default()
	defer p.Release()
	srv := NewGServer("tcp", ":6379", true, p)

	t.Fatal(gnet.Run(srv, srv.net+"://"+srv.addr, gnet.WithMulticore(srv.multicore)))
}
