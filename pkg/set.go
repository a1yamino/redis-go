package pkg

import (
	"strconv"
	"sync"
)

type Set struct {
	m map[string]struct{}
}

func SAddHandler(w IWriter, args []Value) bool {
	if len(args) < 2 {
		w.WriteError("ERR wrong number of arguments for 'sadd' command")
		return false
	}

	key := args[0].String()
	values := args[1:]

	dbMu.Lock()
	e, ok := db[key]
	if ok {
		if e.typ != _Set {
			w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
		e.Lock()
		set := e.value.(*Set)
		for _, v := range values {
			set.m[v.String()] = struct{}{}
		}
		e.Unlock()
	} else {
		s := &Set{make(map[string]struct{})}
		for _, v := range values {
			s.m[v.String()] = struct{}{}
		}
		db[key] = &entry{_Set, s, sync.RWMutex{}}
	}
	dbMu.Unlock()

	w.WriteInteger(len(values))
	return true
}

func SCardHandler(w IWriter, args []Value) bool {
	if len(args) != 1 {
		w.WriteError("ERR wrong number of arguments for 'scard' command")
		return false
	}

	key := args[0].String()

	dbMu.RLock()
	e, ok := db[key]
	if !ok {
		w.WriteInteger(0)
		dbMu.RUnlock()
		return true
	}
	if e.typ != _Set {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.RUnlock()
		return false
	}
	e.RLock()
	set := e.value.(*Set)
	w.WriteInteger(len(set.m))
	e.RUnlock()
	dbMu.RUnlock()
	return true
}

func SIsMemberHandler(w IWriter, args []Value) bool {
	if len(args) != 2 {
		w.WriteError("ERR wrong number of arguments for 'sismember' command")
		return false
	}

	key := args[0].String()
	member := args[1].String()

	dbMu.RLock()
	e, ok := db[key]
	if !ok {
		w.WriteInteger(0)
		dbMu.RUnlock()
		return true
	}
	if e.typ != _Set {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.RUnlock()
		return false
	}
	e.RLock()
	set := e.value.(*Set)
	_, ok = set.m[member]
	w.WriteInteger(bool2int(ok))
	e.RUnlock()
	dbMu.RUnlock()
	return true
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func SMembersHandler(w IWriter, args []Value) bool {
	if len(args) != 1 {
		w.WriteError("ERR wrong number of arguments for 'smembers' command")
		return false
	}

	key := args[0].String()

	dbMu.RLock()
	e, ok := db[key]
	if !ok {
		w.WriteArray(Value{typ: ARRAY, array: []Value{}})
		dbMu.RUnlock()
		return true
	}
	if e.typ != _Set {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.RUnlock()
		return false
	}
	e.RLock()
	set := e.value.(*Set)
	values := make([]Value, 0, len(set.m))
	for v := range set.m {
		values = append(values, BulkString(v))
	}
	e.RUnlock()
	dbMu.RUnlock()
	w.WriteArray(Value{typ: ARRAY, array: values})
	return true
}

func SRandMemberHandler(w IWriter, args []Value) bool {
	if len(args) < 1 || len(args) > 2 {
		w.WriteError("ERR wrong number of arguments for 'srandmember' command")
		return false
	}

	key := args[0].String()
	count := 1
	if len(args) == 2 {
		var err error
		count, err = strconv.Atoi(args[1].String())
		if err != nil {
			w.WriteError("ERR value is not an integer or out of range")
			return false
		}
	}

	dbMu.RLock()
	e, ok := db[key]
	if !ok {
		w.WriteArray(Value{typ: ARRAY, array: []Value{}})
		dbMu.RUnlock()
		return true
	}
	if e.typ != _Set {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.RUnlock()
		return false
	}
	e.RLock()
	set := e.value.(*Set)
	values := make([]Value, 0, count)
	// TODO: negative count
	for v := range set.m {
		values = append(values, BulkString(v))
		count--
		if count == 0 {
			break
		}
	}
	e.RUnlock()
	dbMu.RUnlock()
	w.WriteArray(Value{typ: ARRAY, array: values})
	return true
}

func SPopHandler(w IWriter, args []Value) bool {
	if len(args) != 1 {
		w.WriteError("ERR wrong number of arguments for 'spop' command")
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
	if e.typ != _Set {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.Unlock()
		return false
	}
	e.Lock()
	set := e.value.(*Set)
	for v := range set.m {
		delete(set.m, v)
		w.WriteBulkString(v)
		break
	}
	if len(set.m) == 0 {
		delete(db, key)
	}
	e.Unlock()
	dbMu.Unlock()
	return true
}

func SRemHandler(w IWriter, args []Value) bool {
	if len(args) < 2 {
		w.WriteError("ERR wrong number of arguments for 'srem' command")
		return false
	}

	key := args[0].String()
	members := args[1:]

	dbMu.Lock()
	e, ok := db[key]
	if !ok {
		w.WriteInteger(0)
		dbMu.Unlock()
		return true
	}
	if e.typ != _Set {
		w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
		dbMu.Unlock()
		return false
	}
	e.Lock()
	set := e.value.(*Set)
	count := 0
	for _, v := range members {
		if _, ok := set.m[v.String()]; ok {
			delete(set.m, v.String())
			count++
		}
	}
	if len(set.m) == 0 {
		delete(db, key)
	}
	e.Unlock()
	dbMu.Unlock()
	w.WriteInteger(count)
	return true
}
