//TODO: (1) Review this memory storage & access

package store

import (
	"errors"
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
	Lock  sync.RWMutex
}

// CreateStorage initializes a new store instance
func CreateStorage() *Store {
	s := &Store{
		Items: make(map[string]Data),
		Lock:  sync.RWMutex{},
	}

	return s
}

// GET gets a value for a key in the store
func (s *Store) GET(key string) (Data, error) {
	s.Lock.RLock()
	defer s.Lock.RUnlock()

	data, ok := s.Items[key]
	if !ok {
		return Data{}, errors.New("Key not found")
	}

	if !data.Expiry.IsZero() && data.Expiry.Before(time.Now()) {
		// run a goroutine for deleting expired key also, it cannot be the default zero
		go s.DEL(key)
		return Data{}, errors.New("Expiration time has passed")
	}

	return data, nil
}

// SET gets a key-value pair and adds it to the storage
func (s *Store) SET(key string, value Data) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	s.Items[key] = value
}

// DEL gets a key and deletes it from storage
func (s *Store) DEL(key string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if data, ok := s.Items[key]; ok {
		if !data.Expiry.IsZero() && data.Expiry.Before(time.Now()) {
			delete(s.Items, key)
			return errors.New("Expiration time has passed")
		}
		delete(s.Items, key)
		return nil
	}

	return errors.New("Key not found")
}
