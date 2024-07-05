package pkg

import (
	"strconv"
	"sync"
)

type ZSet struct {
	zsl  *skipList
	dict map[string]float64
}

func NewZSet() *ZSet {
	return &ZSet{
		zsl:  newSkipList(),
		dict: make(map[string]float64),
	}
}

func ZAddHandler(w IWriter, args []Value) bool {
	if len(args) < 3 {
		w.WriteError("ERR wrong number of arguments for 'zadd' command")
		return false
	}

	key := args[0].String()
	dbMu.Lock()
	e, ok := db[key]
	if ok {
		if e.typ != _ZSet {
			w.WriteError("WRONGTYPE Operation against a key holding the wrong kind of value")
			dbMu.Unlock()
			return false
		}
		e.Lock()
		zset := e.value.(*ZSet)
		for i := 1; i < len(args); i += 2 {
			score, err := strconv.ParseFloat(args[i].String(), 64)
			if err != nil {
				w.WriteError("ERR value is not a valid float")
				e.Unlock()
				dbMu.Unlock()
				return false
			}
			member := args[i+1].String()
			if _, ok := zset.dict[member]; ok {
				zset.zsl.Delete(zset.dict[member], member)
			}
			zset.dict[member] = score
			zset.zsl.Insert(score, member)
		}
		e.Unlock()
	} else {
		zset := NewZSet()
		for i := 1; i < len(args); i += 2 {
			score, err := strconv.ParseFloat(args[i].String(), 64)
			if err != nil {
				w.WriteError("ERR value is not a valid float")
				dbMu.Unlock()
				return false
			}
			member := args[i+1].String()
			zset.dict[member] = score
			zset.zsl.Insert(score, member)
		}
		db[key] = &entry{_ZSet, zset, sync.RWMutex{}}
	}
	dbMu.Unlock()

	w.WriteInteger(len(args) / 2)
	return true
}
