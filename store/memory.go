//TODO: (1) Review this memory storage & access

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
		v := data.Value.(resp.BulkString).Value
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, errors.New("Value is not an integer")
		}
		val++
		sval := strconv.Itoa(val)
		newBS := resp.BulkString{Value: sval, Length: len(sval)}
		data.Value = newBS
		s.Items[key] = data

		return resp.Integer{Value: val}, nil
	}

	return nil, errors.New("Key not found")
}