package pkg

import "sync"

type entry struct {
	typ entryType
	// key   string
	value interface{}
	sync.RWMutex
}

type entryType uint8

const (
	_String entryType = iota
	_List
	_Hash
	_Set
	_ZSet
)

var (
	dbMu sync.RWMutex
	db   = make(map[string]*entry)
)
