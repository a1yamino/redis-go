package pkg

import (
	"sync"
)

var (
	str   = make(map[string]string)
	strMu sync.RWMutex
)

func SetHandler(w IWriter, args []Value) bool {
	if len(args) != 2 {
		w.WriteError("ERR wrong number of arguments for 'set' command")
		return false
	}

	key := args[0].String()
	value := args[1].String()

	dbMu.Lock()
	if _, ok := db[key]; ok {
		if db[key].typ != _String {
			w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
	}
	db[key] = &entry{_String, value, sync.RWMutex{}}
	dbMu.Unlock()

	w.WriteSimpleString("OK")
	return true
}

func GetHandler(w IWriter, args []Value) bool {
	if len(args) != 1 {
		w.WriteError("ERR wrong number of arguments for 'get' command")
		return false
	}

	key := args[0].String()

	dbMu.RLock()
	e, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		w.WriteNull()
		return true
	}

	if e.typ != _String {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		return false
	}

	w.WriteBulkString(e.value.(string))
	return true
}

func DelHandler(w IWriter, args []Value) bool {
	if len(args) != 1 {
		w.WriteError("ERR wrong number of arguments for 'del' command")
		return false
	}

	key := args[0].String()

	dbMu.Lock()
	if _, ok := db[key]; !ok {
		w.WriteInteger(0)
		dbMu.Unlock()
		return true
	}
	delete(db, key)
	dbMu.Unlock()
	w.WriteInteger(1)
	return true
}

func ExistsHandler(w IWriter, args []Value) bool {
	if len(args) < 1 {
		w.WriteError("ERR wrong number of arguments for 'exist' command")
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
	w.WriteInteger(result)
	return true
}
