package main

import (
	"log"
	"redis/myredis"
)

func main() {
	log.Println("listening on :6379...")
	server := myredis.NewServer(&myredis.Config{
		EnableAof: false,
	})
	server.ListenAndServe(":6379")
}
