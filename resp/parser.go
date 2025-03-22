//TODO: (1) Change the return type to be based on the input, like if the input is -me\r\n don't return an array

package resp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// GetCommands extracts all commands from the input string and returns them as a slice of strings
func GetCommands(buf []byte) []string {
	val := string(buf)
	val = strings.ToUpper(val)

	inputs := strings.Split(val, "\r\n")

	fmt.Println("input array: ", inputs)

	return inputs
}

// Deserialize takes a command as input and deserializes it
func Deserialize(cmds string) ([]RESPType, error) {
	// Empty command
	command := strings.Trim(cmds, " ")
	if len(command) == 0 {
		return nil, errors.New("Empty input")
	}

	// No CRLF
	if !strings.HasSuffix(command, "\r\n") {
		return nil, errors.New("No CRLF")
	}

	// Remove the last "\r\n" to remove an extra "" in the inputs array
	command = strings.TrimSuffix(command, "\r\n")

	// No commands
	inputs := strings.Split(command, "\r\n")
	if len(inputs) == 0 {
		return nil, errors.New("Invalid input")
	}

	// returned data
	var data []RESPType

	for i := 0; i < len(inputs); i++ {
		v := inputs[i]

		switch v[0] {
		case '+':
			val, err := DeserializeSimpleString(v)
			if err != nil {
				return nil, err
			}
			data = append(data, val)
		case '-':
			val, err := DeserializeSimpleError(v)
			if err != nil {
				return nil, err
			}
			data = append(data, val)
		case ':':
			val, err := DeserializeInteger(v)
			if err != nil {
				return nil, err
			}
			data = append(data, val)
		case '$':
			val, el, err := DeserializeBulkString(inputs[i:])
			if err != nil {
				return nil, err
			}
			data = append(data, val)
			i += el
		default:
			return nil, errors.New("Invalid datatype")
		}
	}

	return data, nil
}

// DeserializeSimpleString returns a simpleString with value stored in it
func DeserializeSimpleString(command string) (SimpleString, error) {
	if command[0] != '+' {
		return SimpleString{}, errors.New("Not SimpleString datatype")
	}
	s := SimpleString{Value: command[1:]}
	return s, nil
}

// DeserializeSimpleError returns a simpleError with value stored in it
func DeserializeSimpleError(command string) (SimpleError, error) {
	if command[0] != '-' {
		return SimpleError{}, errors.New("Not SimpleError datatype")
	}
	s := SimpleError{Value: command[1:]}
	return s, nil
}

// DeserializeInteger returns a simpleError with value stored in it
func DeserializeInteger(command string) (Integer, error) {
	if command[0] != ':' {
		return Integer{}, errors.New("Not Integer datatype")
	}
	val := command[1:]
	num, err := strconv.Atoi(val)
	if err != nil {
		return Integer{}, err
	}
	s := Integer{Value: num}
	return s, nil
}

// DeserializeBulkString returns a simpleError with value stored in it
func DeserializeBulkString(commands []string) (BulkString, int, error) {
	if commands[0][0] != '$' {
		return BulkString{}, 0, errors.New("Not a BulkString datatype")
	}
	val := commands[0][1:]
	bsLen, err := strconv.Atoi(val)
	if err != nil {
		return BulkString{}, 0, err
	}

	bsLenCopy := bsLen

	if bsLen == 0 {
		return BulkString{Value: ""}, 1, nil
	}

	str := ""
	elementsUsed := 0

	for _, v := range commands[1:] {
		if len(v) > bsLen {
			return BulkString{}, 0, errors.New("Invalid Bulk String: length mismatch")
		} else if len(v) == bsLen {
			str += v
			elementsUsed++
			bsLen = 0
			break
		} else {
			str += v
			elementsUsed++
			str += "\r"
			bsLen -= len(str)

			if bsLen != 0 {
				str += "\n"
				bsLen--
				if bsLen == 0 {
					break
				}
			}
		}
	}

	return BulkString{Value: str, Length: bsLenCopy}, elementsUsed, nil
}
