package store

import (
	"strconv"
	"testing"
	"time"

	"github.com/DNahar74/my-redis/resp"
)

func TestGetandSet(t *testing.T) {
	s := CreateStorage()

	key := "testKey"
	value := resp.BulkString{Value: "testValue", Length: len("testValue")}

	data := Data{Value: value}

	s.SET(key, data)

	result, err := s.GET(key)
	if err != nil {
		t.Errorf("Error getting key: %v", err)
	}

	if result.Value != value {
		t.Errorf("Expected value %v, but got %v", value, result.Value)
	}
}

func TestExpiry(t *testing.T) {
	s := CreateStorage()

	key := "testKey"
	value := resp.BulkString{Value: "testValue", Length: len("testValue")}
	data := Data{
		Value:  value,
		Expiry: time.Now().Add(1 * time.Second),
	}

	s.SET(key, data)

	time.Sleep(2 * time.Second)

	_, err := s.GET(key)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestBulkExpiry(t *testing.T) {
	s := CreateStorage()
	numkeys := 100000

	for i := range numkeys {
		key := "key" + strconv.Itoa(i+1)
		val := "value" + strconv.Itoa(i+1)
		value := resp.BulkString{Value: val, Length: len(val)}

		data := Data{
			Value:  value,
			Expiry: time.Now().Add(1 * time.Second),
		}

		s.SET(key, data)
	}

	time.Sleep(2 * time.Second)

	for i := range numkeys {
		key := "key" + strconv.Itoa(i+1)
		_, err := s.GET(key)
		if err == nil {
			t.Errorf("Expected key %v to be expired", key)
		}
	}
}

func TestDelete(t *testing.T) {
	s := CreateStorage()

	key := "testKey"
	value := resp.BulkString{Value: "testValue", Length: len("testValue")}

	data := Data{Value: value}

	s.SET(key, data)

	err := s.DEL(key)
	if err != nil {
		t.Errorf("There was an error deleting the key: %v", err)
	}
}

func TestGetafterDelete(t *testing.T) {
	s := CreateStorage()

	key := "testKey"
	value := resp.BulkString{Value: "testValue", Length: len("testValue")}

	data := Data{Value: value}

	s.SET(key, data)

	err := s.DEL(key)
	if err != nil {
		t.Errorf("There was an error deleting the key: %v", err)
	}

	_, err = s.GET(key)
	if err == nil {
		t.Errorf("Expected key %v to be deleted", key)
	}
}

func TestBulkDeleteExpiry(t *testing.T) {
	s := CreateStorage()
	numkeys := 100000

	for i := range numkeys {
		key := "key" + strconv.Itoa(i+1)
		val := "value" + strconv.Itoa(i+1)
		value := resp.BulkString{Value: val, Length: len(val)}

		data := Data{
			Value:  value,
			Expiry: time.Now().Add(1 * time.Second),
		}

		s.SET(key, data)
	}

	time.Sleep(2 * time.Second)

	for i := range numkeys {
		key := "key" + strconv.Itoa(i+1)
		err := s.DEL(key)
		if err == nil {
			t.Errorf("Expected key %v to be expired", key)
		}
	}
}