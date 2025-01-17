package pkg

import (
	. "container/list"
	"strconv"
	"sync"
)

type list = *List

func LPushHandler(w IWriter, args []Value) bool {
	if len(args) < 2 {
		w.WriteError("ERR wrong number of arguments for 'lpush' command")
		return false
	}

	key := args[0].String()
	values := args[1:]

	dbMu.Lock()
	e, ok := db[key]
	if ok {
		if e.typ != _List {
			w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
		e.Lock()
		lst := e.value.(*qlist)
		for _, v := range values {
			lst.pushLeft(v.String())
		}
		e.Unlock()
	} else {
		l := &qlist{}
		for _, v := range values {
			l.pushLeft(v.String())
		}
		db[key] = &entry{_List, l, sync.RWMutex{}}
	}
	dbMu.Unlock()

	w.WriteInteger(len(values))
	return true
}

func RPushHandler(w IWriter, args []Value) bool {
	if len(args) < 2 {
		w.WriteError("ERR wrong number of arguments for 'rpush' command")
		return false
	}

	key := args[0].String()
	values := args[1:]

	dbMu.Lock()
	e, ok := db[key]
	if ok {
		if e.typ != _List {
			w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
		e.Lock()
		lst := e.value.(*qlist)
		for _, v := range values {
			lst.pushRight(v.String())
		}
		e.Unlock()
	} else {
		l := &qlist{}
		for _, v := range values {
			l.pushRight(v.String())
		}
		db[key] = &entry{_List, l, sync.RWMutex{}}
	}
	dbMu.Unlock()

	w.WriteInteger(len(values))
	return true
}

func LPopHandler(w IWriter, args []Value) bool {
	if len(args) != 1 {
		w.WriteError("ERR wrong number of arguments for 'lpop' command")
		return false
	}

	key := args[0].String()

	dbMu.Lock()
	e, ok := db[key]
	if !ok {
		w.WriteNull()
		dbMu.Unlock()
		return true
	}
	if e.typ != _List {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.Unlock()
		return false
	}
	e.Lock()
	lst := e.value.(*qlist)
	if lst.len == 0 {
		delete(db, key)
		w.WriteNull()
		e.Unlock()
		dbMu.Unlock()
		return true
	}
	v := lst.popLeft()
	e.Unlock()
	dbMu.Unlock()

	w.WriteBulkString(v)
	return true
}

func RPopHandler(w IWriter, args []Value) bool {
	if len(args) != 1 {
		w.WriteError("ERR wrong number of arguments for 'rpop' command")
		return false
	}

	key := args[0].String()

	dbMu.Lock()
	e, ok := db[key]
	if !ok {
		w.WriteNull()
		dbMu.Unlock()
		return true
	}
	if e.typ != _List {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.Unlock()
		return false
	}
	e.Lock()
	lst := e.value.(*qlist)
	if lst.len == 0 {
		delete(db, key)
		w.WriteNull()
		e.Unlock()
		dbMu.Unlock()
		return true
	}
	v := lst.popRight()
	e.Unlock()
	dbMu.Unlock()

	w.WriteBulkString(v)
	return true
}

func LLenHandler(w IWriter, args []Value) bool {
	if len(args) != 1 {
		w.WriteError("ERR wrong number of arguments for 'llen' command")
		return false
	}

	key := args[0].String()

	dbMu.RLock()
	e, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		w.WriteInteger(0)
		return true
	}
	if e.typ != _List {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		return false
	}

	e.RLock()
	lst := e.value.(*qlist)
	l := lst.len
	e.RUnlock()

	w.WriteInteger(l)
	return true
}

func LRangeHandler(w IWriter, args []Value) bool {
	if len(args) != 3 {
		w.WriteError("ERR wrong number of arguments for 'lrange' command")
		return false
	}

	key := args[0].String()
	start, err := strconv.Atoi(args[1].String())
	if err != nil {
		w.WriteError("ERR value is not an integer or out of range")
		return false
	}
	stop, err := strconv.Atoi(args[2].String())
	if err != nil {
		w.WriteError("ERR value is not an integer or out of range")
		return false
	}

	dbMu.RLock()
	e, ok := db[key]
	dbMu.RUnlock()

	if !ok {
		w.WriteArray(Value{typ: ARRAY, array: []Value{}})
		return true
	}
	if e.typ != _List {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		return false
	}

	e.RLock()
	lst := e.value.(*qlist)
	l := lst.len
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
		w.WriteArray(Value{typ: ARRAY, array: []Value{}})
		return true
	}
	if stop >= l {
		stop = l - 1
	}

	values := make([]Value, 0, stop-start+1)
	qrs := lst.getRange(start, stop)
	e.RUnlock()

	for _, qr := range qrs {
		if qr.direction == Left {
			for i := len(qr.data) - 1; i >= 0; i-- {
				values = append(values, BulkString(qr.data[i]))
			}
		} else {
			for _, v := range qr.data {
				values = append(values, BulkString(v))
			}

		}
	}

	w.WriteArray(Value{typ: ARRAY, array: values})
	return true
}

func LTrimHandler(w IWriter, args []Value) bool {
	if len(args) != 3 {
		w.WriteError("ERR wrong number of arguments for 'ltrim' command")
		return false
	}

	key := args[0].String()
	start, err := strconv.Atoi(args[1].String())
	if err != nil {
		w.WriteError("ERR value is not an integer or out of range")
		return false
	}
	stop, err := strconv.Atoi(args[2].String())
	if err != nil {
		w.WriteError("ERR value is not an integer or out of range")
		return false
	}

	dbMu.Lock()
	e, ok := db[key]
	if !ok {
		w.WriteSimpleString("OK")
		dbMu.Unlock()
		return true
	}
	if e.typ != _List {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
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
		w.WriteSimpleString("OK")
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

	w.WriteSimpleString("OK")
	return true
}
