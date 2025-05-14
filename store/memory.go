//TODO: (1) Review this memory storage & access

package store

import (
	"sync"
	"time"

	"github.com/DNahar74/my-redis/resp"
)

// Data represents a key-value pair in the store
type Data struct {
	Value  resp.RESPType
	Expiry time.Time
}

// Store is a map of keys to Data items
type Store struct {
	Items map[string]Data
	Lock sync.RWMutex
}

// CreateStorage initializes a new store instance
func CreateStorage() *Store {
	s := &Store{
		Items: make(map[string]Data),
		Lock: sync.RWMutex{},
	}

	return s
}