//* ALWAYS DO THE LOCK BEFORE THE CHECK (TOCTOU BUG)

package command

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/DNahar74/PulseDB/internal/resp"
	"github.com/DNahar74/PulseDB/internal/store"
)

func handlePING() resp.Type {
	return resp.SimpleString{Value: "PONG"}
}

func handleECHO(value resp.Type) (resp.Type, error) {
	if cmd, ok := value.(resp.BulkString); ok {
		return resp.BulkString{Value: cmd.Value, Length: len(cmd.Value)}, nil
	}

	return nil, errors.New("invalid input type")
}

func handleSET(k, v resp.Type, specs ...resp.Type) (resp.Type, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("invalid key type")
	}

	_, err := strconv.Atoi(key.Value)
	if err == nil {
		return nil, errors.New("key cannot be a number")
	}

	value, ok := v.(resp.BulkString)
	if !ok {
		return nil, errors.New("invalid value type")
	}

	storageData := store.Data{Value: value}

	val, err := strconv.Atoi(value.Value)
	if err == nil {
		storageData.Value = resp.Integer{Value: val}
	}

	specifics := make([]resp.BulkString, 0)

	for _, val := range specs {
		s, ok := val.(resp.BulkString)
		if !ok {
			return nil, errors.New("invalid specs type")
		}

		specifics = append(specifics, s)
	}

	for i := 0; i < len(specifics); i += 2 {
		switch specifics[i].Value {
		case "EX":
			val, err := strconv.Atoi(specifics[i+1].Value)
			if err != nil {
				return nil, err
			}

			storageData.Expiry = time.Now().Add(time.Duration(val) * time.Second)
		default:
			return nil, errors.New("unknown specifier")
		}
	}

	fmt.Println(storageData)

	redisStore.SET(key.Value, storageData)

	return resp.SimpleString{Value: "OK"}, nil
}

func handleGET(k resp.Type) (resp.Type, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("invalid key type")
	}

	data, err := redisStore.GET(key.Value)
	if err != nil {
		return nil, err
	}

	return data.Value, nil
}

func handleDEL(k resp.Type) (resp.Type, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("invalid key type")
	}

	err := redisStore.DEL(key.Value)
	if err != nil {
		return nil, err
	}

	return resp.SimpleString{Value: "OK"}, nil
}

func handleIncr(k resp.Type) (resp.Type, error) {
	key, ok := k.(resp.BulkString)
	if !ok {
		return nil, errors.New("invalid key type")
	}

	val, err := redisStore.INCR(key.Value)
	if err != nil {
		return nil, err
	}

	return val, nil
}
