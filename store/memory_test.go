package store

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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

func TestConcurrentAccess(t *testing.T) {
	s := CreateStorage()
	key := "concurrentKey"
	value := resp.BulkString{Value: "value", Length: 5}
	data := Data{Value: value}

	s.SET(key, data)

	var wg sync.WaitGroup
	for range 100000 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.GET(key)
		}()
	}
	wg.Wait()
}

func TestOverwriteValues(t *testing.T) {
	s := CreateStorage()

	key := "overwrite key"
	value1 := resp.BulkString{Value: "Hello1", Length: 6}
	value2 := resp.BulkString{Value: "Hello2", Length: 6}

	data1 := Data{Value: value1}
	data2 := Data{Value: value2}

	s.SET(key, data1)
	s.SET(key, data2)

	data, err := s.GET(key)
	if err != nil {
		t.Errorf("Error getting key %v", err)
	}

	if data.Value != data2.Value {
		t.Errorf("Incorrect value for key %v", key)
	}
}

func BenchmarkGetSet(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	for range b.N {
		s.SET(key, data)
		_, err := s.GET(key)
		if err != nil {
			b.Errorf("Error in getting key %v", key)
		}
	}
}

func BenchmarkGetSetDynamicKeys(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	for i := range b.N {
		s.SET(key+strconv.Itoa(i), data)
		_, err := s.GET(key + strconv.Itoa(i))
		if err != nil {
			b.Errorf("Error in getting key %v", key+strconv.Itoa(i))
		}
	}
}

func BenchmarkGet(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	s.SET(key, data)

	for range b.N {
		_, err := s.GET(key)
		if err != nil {
			b.Errorf("Error in getting key %v", key)
		}
	}
}

func BenchmarkSet(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	for range b.N {
		s.SET(key, data)
	}
}

func BenchmarkSetDynamicKeys(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	for i := range b.N {
		s.SET(key+strconv.Itoa(i), data)
	}
}

func BenchmarkSetSameKeyParallel(b *testing.B) {
	s := CreateStorage()
	key := "trialkey"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.SET(key, data)
		}
	})
}

func BenchmarkSetDynamicKeysParallel(b *testing.B) {
	s := CreateStorage()
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	var counter uint64

	//? Always create locally changing variables inside goroutines
	b.RunParallel(func(pb *testing.PB) {
		var sb strings.Builder

		for pb.Next() {
			sb.Reset() // Clear the previous string built
			sb.WriteString("benchmarkKey")
			i := atomic.AddUint64(&counter, 1)

			str := strconv.FormatUint(i, 10)
			sb.WriteString(str)
			s.SET(sb.String(), data)
		}
	})
}
