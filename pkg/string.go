package pkg

import (
	"sync"
)

var (
	str   = make(map[string]string)
	strMu sync.RWMutex
)

func SetHandler(conn *Conn, args []Value) bool {
	if len(args) != 2 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'set' command")
		return false
	}

	key := args[0].String()
	value := args[1].String()

	dbMu.Lock()
	if _, ok := db[key]; ok {
		if db[key].typ != _String {
			conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
	}
	db[key] = &entry{_String, value, sync.RWMutex{}}
	dbMu.Unlock()

	conn.Writer.WriteSimpleString("OK")
	return true
}

func GetHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'get' command")
		return false
	}

	key := args[0].String()

	dbMu.RLock()
	e, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		conn.Writer.WriteNull()
		return true
	}

	if e.typ != _String {
		conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		return false
	}

	conn.Writer.WriteBulkString(e.value.(string))
	return true
}

func DelHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'del' command")
		return false
	}

	key := args[0].String()

	dbMu.Lock()
	if _, ok := db[key]; !ok {
		conn.Writer.WriteInteger(0)
		dbMu.Unlock()
		return true
	}
	delete(db, key)
	dbMu.Unlock()
	conn.Writer.WriteInteger(1)
	return true
}

func ExistsHandler(conn *Conn, args []Value) bool {
	if len(args) < 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'exist' command")
		return false
	}

	result := 0

	for _, arg := range args {
		key := arg.String()

		dbMu.RLock()
		_, ok := db[key]
		dbMu.RUnlock()
		if ok {
			result++
		}
	}
	conn.Writer.WriteInteger(result)
	return true
}
