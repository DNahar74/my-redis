// TODO: (1) Check for race condition in delete command

package command

import (
	"errors"

	"github.com/DNahar74/my-redis/resp"
	"github.com/DNahar74/my-redis/store"
)

func handlePING() resp.RESPType {
	return resp.SimpleString{Value: "PONG"}
}

func handleECHO(value resp.RESPType) (resp.RESPType, error) {
	if cmd, ok := value.(resp.BulkString); ok {
		return resp.BulkString{Value: cmd.Value, Length: len(cmd.Value)}, nil
	}

	return nil, errors.New("Invalid input type")
}

func handleSET(k, v resp.RESPType) (resp.RESPType, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid key type")
	}

	value, ok := v.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid value type")
	}

	//? Add a write lock to prevent reads/writes during a write
	redisStore.Lock.Lock()
	defer redisStore.Lock.Unlock()

	redisStore.Items[key.Value] = store.Data{Value: value}
	// redisStore.Items[key.Value] = store.Data{Value: value, Expiry: time.Now()}

	return resp.SimpleString{Value: "OK"}, nil
}

func handleGET(k resp.RESPType) (resp.RESPType, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid key type")
	}

	//? Add a read lock to prevent writes during a read
	//? It blocks writers but readers can proceed
	redisStore.Lock.RLock()
	defer redisStore.Lock.RUnlock()

	if data, ok := redisStore.Items[key.Value]; ok {
		// if data.Expiry.Before(time.Now()) {
    //   delete(redisStore.Items, key.Value)
    //   return nil, errors.New("Key expired")
    // }
		return data.Value, nil
	}

	return nil, errors.New("Key not found")
}

func handleDEL(k resp.RESPType) (resp.RESPType, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid key type")
	}

	if _, ok := redisStore.Items[key.Value]; ok {
		//? Add a write lock to prevent reads/writes during a delete
		//? Also, the existance of the key is checked before deleting (maybe an issue, if there is write then delete, but the order changes due to a race condition)
		redisStore.Lock.Lock()
		defer redisStore.Lock.Unlock()

		delete(redisStore.Items, key.Value)
		return resp.SimpleString{Value: "OK"}, nil
	}

	return nil, errors.New("Key not found")
}