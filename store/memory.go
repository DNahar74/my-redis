//TODO: (1) Review this memory storage & access
//TODO: (2) Change all the data storage from copies to pointer references. make Data.value a pointer to resp.RESPType

package store

import (
	"errors"
	"strconv"
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
	Items   map[string]Data
	Lock    sync.RWMutex
	AOFChan chan string
}

// CreateStorage initializes a new store instance
func CreateStorage() *Store {
	s := &Store{
		Items: make(map[string]Data),
		Lock:  sync.RWMutex{},
		AOFChan: make(chan string, 100000),	// 100000 ops/sec
	}

	return s
}

// GET gets a value for a key in the store
func (s *Store) GET(key string) (Data, error) {
	s.Lock.RLock()

	data, ok := s.Items[key]
	if !ok {
		s.Lock.RUnlock()
		return Data{}, errors.New("Key not found")
	}

	if !data.Expiry.IsZero() && data.Expiry.Before(time.Now()) {
		// run a goroutine for deleting expired key also, it cannot be the default zero

		//? Running a DEL goroutine has issues because it tries to upgrade a read lock to a write lock which is not safe or predictable
		//? To prevent this release the read-lock, create a w-lock & delete
		// go s.DEL(key)

		s.Lock.RUnlock()

		s.Lock.Lock()
		delete(s.Items, key)
		s.Lock.Unlock()

		return Data{}, errors.New("Expiration time has passed")
	}

	s.Lock.RUnlock()

	if iv, ok := data.Value.(resp.Integer); ok {
		v := strconv.Itoa(iv.Value)
		return Data{
			Value:  resp.BulkString{Value: v, Length: len(v)},
			Expiry: data.Expiry,
		}, nil
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
		//? The checking & deletion are in this order because it is impossible to check stuff after deletion
		if !data.Expiry.IsZero() && data.Expiry.Before(time.Now()) {
			delete(s.Items, key)
			return errors.New("Expiration time has passed")
		}
		delete(s.Items, key)
		return nil
	}

	return errors.New("Key not found")
}

// INCR increments the value of a key
func (s *Store) INCR(key string) (resp.RESPType, error) {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if data, ok := s.Items[key]; ok {
		if !data.Expiry.IsZero() && data.Expiry.Before(time.Now()) {
			delete(s.Items, key)
			return nil, errors.New("Expiration time has passed")
		}

		if val, ok := data.Value.(resp.Integer); ok {
			val.Value++
			data.Value = val
			s.Items[key] = data
			return val, nil
		}

		return nil, errors.New("Value is not an integer")
	}

	return nil, errors.New("Key not found")
}
