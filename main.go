package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	log.Println("listening on :6379...")
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}

	aof, err := NewAof("database.aof")
	if err != nil {
		log.Fatal(err)
	}
	defer aof.Close()

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

		if value.typ != "array" {
			fmt.Println("invalid command, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("invalid command, expected non-empty array")
			continue
		}

		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Println("unknown command: ", cmd)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if cmd == "SET" || cmd == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
