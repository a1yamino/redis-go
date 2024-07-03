package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"redis/pkg"

	"github.com/panjf2000/gnet/v2"
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
	srv := pkg.NewGServer("tcp", ":6379", true)
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	log.Fatal(gnet.Run(srv, srv.ProtoAddr(), gnet.WithMulticore(srv.Multicore())))
}
