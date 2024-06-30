package main

import (
	"sync"
)

var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
}

func ping(args []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	SETsMu.Lock()
	defer SETsMu.Unlock()

	SETs[args[0].bulk] = args[1].bulk
	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	SETsMu.RLock()
	defer SETsMu.RUnlock()

	v, ok := SETs[args[0].bulk]
	if !ok {
		return Value{typ: "bull"}
	}
	return Value{typ: "bulk", bulk: v}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}
	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	defer HSETsMu.Unlock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}
	hash := args[0].bulk
	key := args[1].bulk

	HSETsMu.RLock()
	defer HSETsMu.RUnlock()
	v, ok := HSETs[hash][key]
	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: v}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}
	hash := args[0].bulk

	HSETsMu.RLock()
	defer HSETsMu.RUnlock()
	var values []Value
	for key, value := range HSETs[hash] {
		values = append(values, Value{typ: "bulk", bulk: key})
		values = append(values, Value{typ: "bulk", bulk: value})
	}
	return Value{typ: "array", array: values}
}
