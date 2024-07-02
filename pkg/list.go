package pkg

import (
	. "container/list"
	"strconv"
	"sync"
)

type list = *List

func LPushHandler(conn *Conn, args []Value) bool {
	if len(args) < 2 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'lpush' command")
		return false
	}

	key := args[0].String()
	values := args[1:]

	dbMu.Lock()
	e, ok := db[key]
	if ok {
		if e.typ != _List {
			conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
		e.Lock()
		lst := e.value.(list)
		for _, v := range values {
			lst.PushFront(v.String())
		}
		e.Unlock()
	} else {
		l := list(New())
		for _, v := range values {
			l.PushFront(v.String())
		}
		db[key] = &entry{_List, l, sync.RWMutex{}}
	}
	dbMu.Unlock()

	conn.Writer.WriteInteger(len(values))
	return true
}

func RPushHandler(conn *Conn, args []Value) bool {
	if len(args) < 2 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'rpush' command")
		return false
	}

	key := args[0].String()
	values := args[1:]

	dbMu.Lock()
	e, ok := db[key]
	if ok {
		if e.typ != _List {
			conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
		e.Lock()
		lst := e.value.(list)
		for _, v := range values {
			lst.PushBack(v.String())
		}
		e.Unlock()
	} else {
		l := list(New())
		for _, v := range values {
			l.PushBack(v.String())
		}
		db[key] = &entry{_List, l, sync.RWMutex{}}
	}
	dbMu.Unlock()

	conn.Writer.WriteInteger(len(values))
	return true
}

func LPopHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'lpop' command")
		return false
	}

	key := args[0].String()

	dbMu.Lock()
	e, ok := db[key]
	if !ok {
		conn.Writer.WriteNull()
		dbMu.Unlock()
		return true
	}
	if e.typ != _List {
		conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.Unlock()
		return false
	}
	e.Lock()
	lst := e.value.(list)
	if lst.Len() == 0 {
		conn.Writer.WriteNull()
		e.Unlock()
		dbMu.Unlock()
		return true
	}
	v := lst.Remove(lst.Front()).(string)
	e.Unlock()
	dbMu.Unlock()

	conn.Writer.WriteBulkString(v)
	return true
}

func RPopHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'rpop' command")
		return false
	}

	key := args[0].String()

	dbMu.Lock()
	e, ok := db[key]
	if !ok {
		conn.Writer.WriteNull()
		dbMu.Unlock()
		return true
	}
	if e.typ != _List {
		conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.Unlock()
		return false
	}
	e.Lock()
	lst := e.value.(list)
	if lst.Len() == 0 {
		conn.Writer.WriteNull()
		e.Unlock()
		dbMu.Unlock()
		return true
	}
	v := lst.Remove(lst.Back()).(string)
	e.Unlock()
	dbMu.Unlock()

	conn.Writer.WriteBulkString(v)
	return true
}

func LLenHandler(conn *Conn, args []Value) bool {
	if len(args) != 1 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'llen' command")
		return false
	}

	key := args[0].String()

	dbMu.RLock()
	e, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		conn.Writer.WriteInteger(0)
		return true
	}
	if e.typ != _List {
		conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		return false
	}

	e.RLock()
	lst := e.value.(list)
	l := lst.Len()
	e.RUnlock()

	conn.Writer.WriteInteger(l)
	return true
}

func LRangeHandler(conn *Conn, args []Value) bool {
	if len(args) != 3 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'lrange' command")
		return false
	}

	key := args[0].String()
	start, err := strconv.Atoi(args[1].String())
	if err != nil {
		conn.Writer.WriteError("ERR value is not an integer or out of range")
		return false
	}
	stop, err := strconv.Atoi(args[2].String())
	if err != nil {
		conn.Writer.WriteError("ERR value is not an integer or out of range")
		return false
	}

	dbMu.RLock()
	e, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		conn.Writer.WriteArray(Value{typ: ARRAY, array: []Value{}})
		return true
	}
	if e.typ != _List {
		conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		return false
	}

	e.RLock()
	lst := e.value.(list)
	l := lst.Len()
	if start < 0 {
		start = l + start
	}
	if stop < 0 {
		stop = l + stop
	}
	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 0
	}
	if start >= l {
		e.RUnlock()
		conn.Writer.WriteArray(Value{typ: ARRAY, array: []Value{}})
		return true
	}
	if stop >= l {
		stop = l - 1
	}

	values := make([]Value, stop-start+1)
	i := 0
	for e := lst.Front(); e != nil; e = e.Next() {
		if i >= start && i <= stop {
			values[i-start] = BulkString(e.Value.(string))
		}
		i++
	}
	e.RUnlock()
	conn.Writer.WriteArray(Value{typ: ARRAY, array: values})
	return true
}

func LTrimHandler(conn *Conn, args []Value) bool {
	if len(args) != 3 {
		conn.Writer.WriteError("ERR wrong number of arguments for 'ltrim' command")
		return false
	}

	key := args[0].String()
	start, err := strconv.Atoi(args[1].String())
	if err != nil {
		conn.Writer.WriteError("ERR value is not an integer or out of range")
		return false
	}
	stop, err := strconv.Atoi(args[2].String())
	if err != nil {
		conn.Writer.WriteError("ERR value is not an integer or out of range")
		return false
	}

	dbMu.Lock()
	e, ok := db[key]
	if !ok {
		conn.Writer.WriteSimpleString("OK")
		dbMu.Unlock()
		return true
	}
	if e.typ != _List {
		conn.Writer.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.Unlock()
		return false
	}

	e.Lock()
	lst := e.value.(list)
	l := lst.Len()
	if start < 0 {
		start = l + start
	}
	if stop < 0 {
		stop = l + stop
	}
	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 0
	}
	if start >= l {
		lst.Init()
		delete(db, key)
		e.Unlock()
		dbMu.Unlock()
		conn.Writer.WriteSimpleString("OK")
		return true
	}
	if stop >= l {
		stop = l - 1
	}

	i := 0
	for e := lst.Front(); e != nil; e = e.Next() {
		if i < start || i > stop {
			lst.Remove(e)
		}
		i++
	}
	e.Unlock()
	dbMu.Unlock()

	conn.Writer.WriteSimpleString("OK")
	return true
}
