package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	log.Println("listening on :6379...")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := l.Accept()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v\n", value)

		conn.Write([]byte("+OK\r\n"))
	}
}
