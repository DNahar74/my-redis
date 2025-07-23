package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DNahar74/PulseDB/internal/resp"
	"github.com/DNahar74/PulseDB/internal/store"
)

var redisStore *store.Store

// InitStore passes the RedisStore global variable's pointer for access in this package
func InitStore(rs *store.Store) {
	redisStore = rs
}

// HandleCommands takes a Type and handles it based on the command type
func HandleCommands(commands resp.Type) (resp.Type, error) {
	switch commands.(type) {
	case resp.SimpleString:
		val, err := handleSimpleString(commands)
		if err != nil {
			return nil, err
		}

		return val, nil
	case resp.SimpleError:
		val, err := handleSimpleError(commands)
		if err != nil {
			return nil, err
		}

		return val, nil
	case resp.Integer:
		val, err := handleInteger(commands)
		if err != nil {
			return nil, err
		}

		return val, nil
	case resp.BulkString:
		val, err := handleBulkString(commands)
		if err != nil {
			return nil, err
		}

		return val, nil
	case resp.Array:
		val, err := handleArray(commands)
		if err != nil {
			return nil, err
		}

		return val, nil
	default:
		fmt.Printf("Unknown type: %T", commands)
	}

	return nil, nil
}

func handleSimpleString(command resp.Type) (resp.Type, error) {
	if str, ok := command.(resp.SimpleString); ok {
		return str, nil
	}
	return nil, errors.New("invalid datatype")
}

func handleSimpleError(command resp.Type) (resp.Type, error) {
	if str, ok := command.(resp.SimpleError); ok {
		return str, nil
	}
	return nil, errors.New("invalid datatype")
}

func handleInteger(command resp.Type) (resp.Type, error) {
	if str, ok := command.(resp.Integer); ok {
		return str, nil
	}
	return nil, errors.New("invalid datatype")
}

func handleBulkString(command resp.Type) (resp.Type, error) {
	if str, ok := command.(resp.BulkString); ok {
		return str, nil
	}
	return nil, errors.New("invalid datatype")
}

func handleArray(command resp.Type) (resp.Type, error) {
	if str, ok := command.(resp.Array); ok {
		// If the command is an array, check if the first element is a BulkString(all command names are BulkStrings)
		if cmd, ok := str.Items[0].(resp.BulkString); ok {
			switch strings.ToUpper(cmd.Value) {
			case "PING":
				v := handlePING()
				return v, nil
			case "ECHO":
				if len(str.Items) < 2 {
					return nil, errors.New("ECHO requires an argument")
				}
				v, err := handleECHO(str.Items[1])
				if err != nil {
					return nil, err
				}
				return v, nil
			case "SET":
				if len(str.Items) < 3 {
					return nil, errors.New("SET requires a key and a value")
				} else if len(str.Items) == 3 {
					v, err := handleSET(str.Items[1], str.Items[2])
					if err != nil {
						return nil, err
					}
					return v, nil
				} else {
					v, err := handleSET(str.Items[1], str.Items[2], str.Items[3:]...)
					if err != nil {
						return nil, err
					}
					return v, nil
				}
			case "GET":
				if len(str.Items) < 2 {
					return nil, errors.New("GET requires a key")
				} else if len(str.Items) > 2 {
					return nil, errors.New("GET supports only one key")
				}
				v, err := handleGET(str.Items[1])
				if err != nil {
					return nil, err
				}
				return v, nil
			case "DEL":
				if len(str.Items) < 2 {
					return nil, errors.New("DEL requires a key")
				} else if len(str.Items) > 2 {
					return nil, errors.New("DEL supports only one key")
				}
				v, err := handleDEL(str.Items[1])
				if err != nil {
					return nil, err
				}
				return v, nil
			case "INCR":
				if len(str.Items) < 2 {
					return nil, errors.New("INCR requires a key")
				} else if len(str.Items) > 2 {
					return nil, errors.New("INCR supports only one key")
				}
				v, err := handleIncr(str.Items[1])
				if err != nil {
					return nil, err
				}
				return v, nil
			default:
				return nil, errors.New("unknown Command")
			}
		}
		return nil, errors.New("not a valid command format")
	}
	return nil, errors.New("invalid datatype")
}
