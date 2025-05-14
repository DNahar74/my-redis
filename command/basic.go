//* ALWAYS DO THE LOCK BEFORE THE CHECK (TOCTOU BUG)

package command

import (
	"errors"
	"fmt"
	"strconv"
	"time"

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

func handleSET(k, v resp.RESPType, specs ...resp.RESPType) (resp.RESPType, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid key type")
	}

	_, err := strconv.Atoi(key.Value)
	if err == nil {
		return nil, errors.New("Key cannot be a number")
	}

	value, ok := v.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid value type")
	}

	storageData := store.Data{Value: value}

	specifics := make([]resp.BulkString, 0)

	for _, val := range specs {
		s, ok := val.(resp.BulkString)
		if !ok {
			return nil, errors.New("Invalid specs type")
		}

		specifics = append(specifics, s)
	}

	for i := 0; i < len(specifics); i += 2 {
		if specifics[i].Value == "EX" {
			val, err := strconv.Atoi(specifics[i+1].Value)
			if err != nil {
				return nil, err
			}

			// fmt.Println(val)
			// fmt.Println(time.Now())

			storageData.Expiry = time.Now().Add(time.Duration(val) * time.Second)
		}
	}

	fmt.Println(storageData)

	redisStore.SET(key.Value, storageData)

	return resp.SimpleString{Value: "OK"}, nil
}

func handleGET(k resp.RESPType) (resp.RESPType, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid key type")
	}

	data, err := redisStore.GET(key.Value)
	if err != nil {
		return nil, err
	}

	return data.Value, nil
}

func handleDEL(k resp.RESPType) (resp.RESPType, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid key type")
	}

	err := redisStore.DEL(key.Value)
	if err != nil {
		return nil, err
	}

	return resp.SimpleString{Value: "OK"}, nil
}
