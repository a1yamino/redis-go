package myredis

import (
	"fmt"
	"sync"
)

var (
	hashMu sync.RWMutex
	hash   = make(map[string]map[string]string)
)

func HSetHandler(conn *Conn, args []Value) bool {
	if len(args) != 3 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hset' command")
		return false
	}

	key := args[0].String()
	field := args[1].String()
	value := args[2].String()

	hashMu.Lock()
	if _, ok := hash[key]; !ok {
		hash[key] = make(map[string]string)
	}
	hash[key][field] = value
	hashMu.Unlock()

	conn.Writer.WriteInteger(1)
	return true
}

func HGetHandler(conn *Conn, args []Value) bool {
	if len(args) != 2 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hget' command")
		return false
	}

	key := args[0].String()
	field := args[1].String()

	hashMu.RLock()
	value, ok := hash[key][field]
	hashMu.RUnlock()

	if !ok {
		conn.Writer.WriteNull()
		return true
	}

	conn.Writer.WriteBulkString(value)
	return true
}

func HGetAllHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hgetall' command")
		return false
	}

	key := args[0].String()

	hashMu.RLock()
	fields := hash[key]
	hashMu.RUnlock()

	values := make([]Value, 0, len(fields)*2)
	for field, value := range fields {
		values = append(values, BulkString(field), BulkString(value))
	}

	err := conn.Writer.WriteArray(Value{typ: ARRAY, array: values})
	if err != nil {
		fmt.Printf("write array failed: %v\n", err)
	}
	return true
}

func HDelHandler(conn *Conn, args []Value) bool {
	if len(args) < 2 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hdel' command")
		return false
	}

	key := args[0].String()
	fields := make([]string, 0, len(args)-1)
	for _, arg := range args[1:] {
		fields = append(fields, arg.String())
	}

	hashMu.Lock()
	if _, ok := hash[key]; !ok {
		hashMu.Unlock()
		conn.Writer.WriteInteger(0)
		return true
	}

	var count int
	for _, field := range fields {
		if _, ok := hash[key][field]; ok {
			delete(hash[key], field)
			count++
		}
	}
	hashMu.Unlock()

	conn.Writer.WriteInteger(count)
	return true
}

func HLenHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hlen' command")
		return false
	}

	key := args[0].String()

	hashMu.RLock()
	fields, ok := hash[key]
	hashMu.RUnlock()

	if !ok {
		conn.Writer.WriteInteger(0)
		return true
	}

	conn.Writer.WriteInteger(len(fields))
	return true
}

func HKeysHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hkeys' command")
		return false
	}

	key := args[0].String()

	hashMu.RLock()
	fields, ok := hash[key]
	hashMu.RUnlock()

	if !ok {
		conn.Writer.WriteNull()
		return true
	}

	values := make([]Value, 0, len(fields))
	for field := range fields {
		values = append(values, BulkString(field))
	}

	err := conn.Writer.WriteArray(Value{typ: ARRAY, array: values})
	if err != nil {
		fmt.Printf("write array failed: %v\n", err)
	}
	return true
}

func HValsHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hvals' command")
		return false
	}

	key := args[0].String()

	hashMu.RLock()
	fields, ok := hash[key]
	hashMu.RUnlock()

	if !ok {
		conn.Writer.WriteNull()
		return true
	}

	values := make([]Value, 0, len(fields))
	for _, value := range fields {
		values = append(values, BulkString(value))
	}

	err := conn.Writer.WriteArray(Value{typ: ARRAY, array: values})
	if err != nil {
		fmt.Printf("write array failed: %v\n", err)
	}
	return true
}
