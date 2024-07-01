package main

import (
	"log"
	"redis/pkg"
)

func main() {
	log.Println("listening on :6379...")
	server := pkg.NewServer(&pkg.Config{
		EnableAof: false,
	})
	server.ListenAndServe(":6379")
}
