package myredis

import "sync"

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

	strMu.Lock()
	str[key] = value
	strMu.Unlock()

	conn.Writer.WriteSimpleString("OK")
	return true
}

func GetHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'get' command")
		return false
	}

	key := args[0].String()

	strMu.RLock()
	value, ok := str[key]
	strMu.RUnlock()

	if !ok {
		conn.Writer.WriteNull()
		return true
	}

	conn.Writer.WriteBulkString(value)
	return true
}

func DelHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'del' command")
		return false
	}

	key := args[0].String()

	strMu.Lock()
	if _, ok := str[key]; !ok {
		strMu.Unlock()

		conn.Writer.WriteInteger(0)
		return true
	}
	delete(str, key)
	strMu.Unlock()

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

		strMu.RLock()
		_, ok := str[key]
		strMu.RUnlock()
		if ok {
			result++
		}
	}
	conn.Writer.WriteInteger(result)
	return true
}
