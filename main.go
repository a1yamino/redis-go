package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"redis/pkg"
)

func main() {
	log.Println("listening on :6379...")
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	server := pkg.NewServer(&pkg.Config{
		EnableAof: false,
	})
	log.Fatal(server.ListenAndServe(":6379"))
}
