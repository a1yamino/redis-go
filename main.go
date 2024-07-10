package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"redis/pkg"

	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/pool/goroutine"
)

// func main() {
// 	log.Println("listening on :6379...")
// 	go func() {
// 		log.Println(http.ListenAndServe("localhost:6060", nil))
// 	}()
// 	server := pkg.NewServer(&pkg.Config{
// 		EnableAof: false,
// 	})
// 	log.Fatal(server.ListenAndServe(":6379"))
// }

func main() {
	p := goroutine.Default()
	defer p.Release()
	srv := pkg.NewGServer("tcp", ":6379", true, p)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	log.Fatal(gnet.Run(srv, srv.ProtoAddr(), gnet.WithMulticore(srv.Multicore()), gnet.WithTCPNoDelay(gnet.TCPDelay)))
}
