package command

import "github.com/DNahar74/my-redis/resp"

func handlePING() resp.RESPType {
	return resp.SimpleString{Value: "PONG"}
}
