package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"redis/pkg"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
)

var (
	netpkg = flag.String("netpkg", "net", "gnet or net")
	netmap = map[string]func(){
		"gnet": func() {
			p := goroutine.Default()
			defer p.Release()
			srv := pkg.NewGServer("tcp", ":6379", true, p)
			log.Fatal(gnet.Run(srv, srv.ProtoAddr(), gnet.WithMulticore(srv.Multicore())))
		},
		"net": func() {
			log.Println("listening on :6379...")
			server := pkg.NewServer(&pkg.Config{
				EnableAof: false,
			})
			log.Fatal(server.ListenAndServe(":6379"))
		},
	}
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	flag.Parse()
	if f, ok := netmap[*netpkg]; ok {
		f()
	} else {
		log.Fatal("invalid netpkg")
	}
}
