package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DNahar74/my-redis/resp"
)

// HandleCommands takes a RESPType and handles it based on the command type
func HandleCommands(commands resp.RESPType) (resp.RESPType, error) {
	switch commands.(type) {
	case resp.SimpleString:
		val, err := handleSimpleString(commands)
		if err != nil {
			return nil, err
		}

		fmt.Println("It's a simple string with value:", val)
		return val, nil
	case resp.SimpleError:
		val, err := handleSimpleError(commands)
		if err != nil {
			return nil, err
		}

		fmt.Println("It's a simple error with value:", val)
		return val, nil
	case resp.Integer:
		val, err := handleInteger(commands)
		if err != nil {
			return nil, err
		}

		fmt.Println("It's an integer with value:", val)
		return val, nil
	case resp.BulkString:
		val, err := handleBulkString(commands)
		if err != nil {
			return nil, err
		}

		fmt.Println("It's a bulk string with value:", val)
		return val, nil
	case resp.Array:
		val, err := handleArray(commands)
		if err != nil {
			return nil, err
		}

		fmt.Println("It's an array with elements:", val)
		return val, nil
	default:
		fmt.Printf("Unknown type: %T", commands)
	}

	return nil, nil
}

func handleSimpleString(command resp.RESPType) (resp.RESPType, error) {
	if str, ok := command.(resp.SimpleString); ok {
		return str, nil
	}
	return nil, errors.New("Invalid datatype")
}

func handleSimpleError(command resp.RESPType) (resp.RESPType, error) {
	if str, ok := command.(resp.SimpleError); ok {
		return str, nil
	}
	return nil, errors.New("Invalid datatype")
}

func handleInteger(command resp.RESPType) (resp.RESPType, error) {
	if str, ok := command.(resp.Integer); ok {
		return str, nil
	}
	return nil, errors.New("Invalid datatype")
}

func handleBulkString(command resp.RESPType) (resp.RESPType, error) {
	if str, ok := command.(resp.BulkString); ok {
		return str, nil
	}
	return nil, errors.New("Invalid datatype")
}

func handleArray(command resp.RESPType) (resp.RESPType, error) {
	if str, ok := command.(resp.Array); ok {

		if cmd, ok := str.Items[0].(resp.BulkString); ok {
			switch strings.ToUpper(cmd.Value) {
			case "PING":
				v := handlePING()
				return v, nil
			case "ECHO":
				v, err := handleECHO(str.Items[1])
				if err != nil {
          return nil, err
        }
        return v, nil
			default:
				return nil, errors.New("Unknown Command")
			}
		}
		return nil, errors.New("Not a valid command format")
	}
	return nil, errors.New("Invalid datatype")
}
