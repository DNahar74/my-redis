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

	redisStore.Items[key.Value] = store.Data{Value: value}

	return resp.SimpleString{Value: "OK"}, nil
}

func handleGET(k resp.RESPType) (resp.RESPType, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("Invalid key type")
	}

	if data, ok := redisStore.Items[key.Value]; ok {
		return data.Value, nil
	}

	return nil, errors.New("Key not found")
}
