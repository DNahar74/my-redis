package command

import (
	"errors"

	"github.com/DNahar74/my-redis/resp"
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
