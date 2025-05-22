package store

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/DNahar74/my-redis/resp"
)

func TestGetAndSet(t *testing.T) {
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
	numKeys := 100000

	for i := range numKeys {
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

	for i := range numKeys {
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

func TestGetAfterDelete(t *testing.T) {
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
	numKeys := 100000

	for i := range numKeys {
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

	for i := range numKeys {
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
			_, err := s.GET(key)
			if err != nil {
				t.Errorf("There was an error getting the key: %v", err)
			}
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

func TestExpiryUpdate(t *testing.T) {
	s := CreateStorage()
	key := "key1"
	data1 := Data{Value: resp.BulkString{Value: "val", Length: 3}, Expiry: time.Now().Add(1 * time.Second)}
	data2 := Data{Value: resp.BulkString{Value: "val", Length: 3}, Expiry: time.Now().Add(5 * time.Second)}

	s.SET(key, data1)
	s.SET(key, data2)

	time.Sleep(2 * time.Second)

	_, err := s.GET(key)
	if err != nil {
		t.Errorf("Key should not have expired yet: %v", err)
	}
}

func TestNoExpiry(t *testing.T) {
	s := CreateStorage()
	key := "keyNoExpiry"
	data := Data{Value: resp.BulkString{Value: "permanent", Length: 9}}

	s.SET(key, data)
	time.Sleep(2 * time.Second)

	_, err := s.GET(key)
	if err != nil {
		t.Errorf("Expected key to persist, but got error: %v", err)
	}
}

func TestGetNonExistentKey(t *testing.T) {
	s := CreateStorage()
	_, err := s.GET("unknownKey")
	if err == nil {
		t.Errorf("Expected error for unknown key")
	}
}

func TestDeleteNonExistentKey(t *testing.T) {
	s := CreateStorage()
	err := s.DEL("missingKey")
	if err == nil {
		t.Errorf("Expected error when deleting non-existent key")
	}
}

func TestEmptyValue(t *testing.T) {
	s := CreateStorage()
	key := "emptyKey"
	data := Data{Value: resp.BulkString{Value: "", Length: 0}}

	s.SET(key, data)
	res, err := s.GET(key)
	if err != nil || res.Value.(resp.BulkString).Value != "" {
		t.Errorf("Expected empty string, got %v", res.Value.(resp.BulkString).Value)
	}
}

func TestConcurrentExpiry(t *testing.T) {
	s := CreateStorage()

	key := "expireKey"
	data := Data{
		Value:  resp.BulkString{Value: "value", Length: 5},
		Expiry: time.Now().Add(1 * time.Second),
	}
	s.SET(key, data)

	var wg sync.WaitGroup

	for range 10000 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			time.Sleep(2 * time.Second)
			_, err := s.GET(key)
			if err == nil {
				t.Errorf("Error was expected, but value is accessible after expiry")
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentSetGetDifferentKeys(t *testing.T) {
	s := CreateStorage()
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i)
			data := Data{Value: resp.BulkString{Value: "val", Length: 3}}
			s.SET(key, data)
			_, _ = s.GET(key)
		}(i)
	}
	wg.Wait()
}

func TestIncrement(t *testing.T) {
	s := CreateStorage()

	key := "incrementKey"
	value := resp.Integer{Value: 11}
	data := Data{Value: value}

	s.SET(key, data)

	v, err := s.INCR(key)
	if err != nil {
		t.Errorf("Got an unexpected error incrementing key : %v", err)
	}

	if v.(resp.Integer).Value != 12 {
		t.Errorf("Did not increment successfully")
	}
}

func TestIncrementNegative(t *testing.T) {
	s := CreateStorage()

	key := "incrementKey"
	value := resp.Integer{Value: -11}
	data := Data{Value: value}

	s.SET(key, data)

	v, err := s.INCR(key)
	if err != nil {
		t.Errorf("Got an unexpected error incrementing key : %v", err)
	}

	if v.(resp.Integer).Value != -10 {
		t.Errorf("Did not increment successfully")
	}
}

func TestIncrementExpiry(t *testing.T) {
	s := CreateStorage()

	key := "incrementKey"
	value := resp.Integer{Value: 11}
	data := Data{
		Value:  value,
		Expiry: time.Now().Add(1 * time.Second),
	}
	s.SET(key, data)

	time.Sleep(2 * time.Second)

	_, err := s.INCR(key)
	if err == nil {
		t.Errorf("Expected error incrementing an expired key, but did not get an error")
	}
}

func TestStringIncrement(t *testing.T) {
	s := CreateStorage()

	key := "incrementKey"
	value := resp.BulkString{Value: "hello", Length: 5}
	data := Data{Value: value}

	s.SET(key, data)

	_, err := s.INCR(key)
	if err == nil {
		t.Errorf("Expected an error incrementing a string value, but did not get an error")
	}
}

func TestIncrementConcurrently(t *testing.T) {
	s := CreateStorage()
	var wg sync.WaitGroup

	startVal := 11

	key := "incrementKey"
	value := resp.Integer{Value: startVal}
	data := Data{Value: value}
	s.SET(key, data)

	intervals := 1000000

	for range intervals {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.INCR(key)
			if err != nil {
				t.Errorf("Got an unexpected error incrementing key : %v", err)
			}
		}()
	}
	wg.Wait()

	v, err := s.INCR(key)
	if err != nil {
		t.Errorf("Got an unexpected error incrementing key : %v", err)
	}

	final := v.(resp.Integer).Value
	expected := startVal + intervals + 1

	if final != expected {
		t.Errorf("Expected %d, got %d", expected, final)
	}
}

func BenchmarkGetSet(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	// Reset timer to exclude setup time
	b.ResetTimer()

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

	// Reset timer to exclude setup time
	b.ResetTimer()

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

	// Reset timer to exclude setup time
	b.ResetTimer()

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

	// Reset timer to exclude setup time
	b.ResetTimer()

	for range b.N {
		s.SET(key, data)
	}
}

func BenchmarkSetDynamicKeys(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	// Reset timer to exclude setup time
	b.ResetTimer()

	for i := range b.N {
		s.SET(key+strconv.Itoa(i), data)
	}
}

func BenchmarkSetSameKeyParallel(b *testing.B) {
	s := CreateStorage()
	key := "trial key"
	value := resp.BulkString{Value: "Hello", Length: 5}
	data := Data{Value: value}

	// Reset timer to exclude setup time
	b.ResetTimer()

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

	// Reset timer to exclude setup time
	b.ResetTimer()

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

func BenchmarkIncrement(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.Integer{Value: 0}
	data := Data{Value: value}
	s.SET(key, data)

	// Reset timer to exclude setup time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := s.INCR(key)
		if err != nil {
			b.Fatalf("Unexpected error during INCR: %v", err)
		}
	}
}

func BenchmarkIncrementParallel(b *testing.B) {
	s := CreateStorage()

	key := "benchmarkKey"
	value := resp.Integer{Value: 0}
	data := Data{Value: value}
	s.SET(key, data)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := s.INCR(key)
			if err != nil {
				b.Fatalf("Unexpected error during parallel INCR: %v", err)
			}
		}
	})
}
