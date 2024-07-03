package main

import (
	"log"
	"redis/pkg"
	"testing"
)

func TestXxx(t *testing.T) {
	server := pkg.NewServer(&pkg.Config{
		EnableAof: false,
	})
	log.Fatal(server.ListenAndServe(":6379"))
}
