package pkg

import (
	"testing"

	"github.com/panjf2000/gnet/v2"
)

func TestGnetServer(t *testing.T) {
	srv := NewGServer("tcp", ":6379", true)

	t.Fatal(gnet.Run(srv, srv.net+"://"+srv.addr, gnet.WithMulticore(srv.multicore)))
}
