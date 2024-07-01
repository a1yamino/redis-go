package myredis

import (
	"fmt"
	"sync"
)

// var (
// 	hashMu sync.RWMutex
// 	hash   = make(map[string]map[string]string)
// )

func HSetHandler(conn *Conn, args []Value) bool {
	if len(args) != 3 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hset' command")
		return false
	}

	key := args[0].String()
	field := args[1].String()
	value := args[2].String()

	dbMu.Lock()
	if _, ok := db[key]; ok {
		if db[key].typ != _Hash {
			conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
		db[key].Lock()
		db[key].value.(hash)[field] = value
		db[key].Unlock()
	} else {
		db[key] = &entry{_Hash, hash{field: value}, sync.RWMutex{}}
	}
	dbMu.Unlock()
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

	dbMu.RLock()
	e, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		conn.Writer.WriteNull()
		return true
	}
	e.RLock()
	if e.typ != _Hash {
		conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		e.RUnlock()
		return false
	}

	value, ok := e.value.(hash)[field]
	e.RUnlock()
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

	dbMu.RLock()
	hashEntry := db[key]
	dbMu.RUnlock()

	hashEntry.RLock()
	if hashEntry.typ != _Hash {
		conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		hashEntry.RUnlock()
		return false
	}
	hashV := hashEntry.value.(hash)
	values := make([]Value, 0, len(hashV)*2)
	for field, value := range hashV {
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

	dbMu.RLock()
	hashEntry, ok := db[key]
	if !ok {
		dbMu.RUnlock()
		conn.Writer.WriteInteger(0)
		return true
	}
	dbMu.RUnlock()
	var count int
	hashEntry.Lock()
	hashV := hashEntry.value.(hash)
	for _, field := range fields {
		if _, ok := hashV[field]; ok {
			delete(hashV, field)
			count++
		}
	}
	hashEntry.Unlock()
	conn.Writer.WriteInteger(count)
	return true
}

func HLenHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hlen' command")
		return false
	}

	key := args[0].String()

	dbMu.RLock()
	hashEntry, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		conn.Writer.WriteInteger(0)
		return true
	}

	hashEntry.RLock()
	conn.Writer.WriteInteger(len(hashEntry.value.(hash)))
	hashEntry.RUnlock()
	return true
}

func HKeysHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'hkeys' command")
		return false
	}

	key := args[0].String()

	dbMu.RLock()
	hashEntry, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		conn.Writer.WriteNull()
		return true
	}

	hashEntry.RLock()
	hashV := hashEntry.value.(hash)
	values := make([]Value, 0, len(hashV))
	for field := range hashV {
		values = append(values, BulkString(field))
	}
	hashEntry.RUnlock()
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

	dbMu.RLock()
	hashEntry, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		conn.Writer.WriteNull()
		return true
	}

	hashEntry.RLock()
	hashV := hashEntry.value.(hash)
	values := make([]Value, 0, len(hashV))
	for _, value := range hashV {
		values = append(values, BulkString(value))
	}
	hashEntry.RUnlock()

	err := conn.Writer.WriteArray(Value{typ: ARRAY, array: values})
	if err != nil {
		fmt.Printf("write array failed: %v\n", err)
	}
	return true
}
