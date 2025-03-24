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
func Deserialize(cmds string) (RESPType, error) {
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

	for i := 0; i < len(inputs); i++ {
		v := inputs[i]

		switch v[0] {
		case '+':
			val, err := DeserializeSimpleString(v)
			if err != nil {
				return nil, err
			}
			return val, nil
		case '-':
			val, err := DeserializeSimpleError(v)
			if err != nil {
				return nil, err
			}
			return val, nil
		case ':':
			val, err := DeserializeInteger(v)
			if err != nil {
				return nil, err
			}
			return val, nil
		case '$':
			val, _, err := DeserializeBulkString(inputs[i:])
			if err != nil {
				return nil, err
			}
			return val, nil
		case '*':
			val, _, err := DeserializeArray(inputs[i:])
			if err != nil {
				return nil, err
			}
			return val, nil
		default:
			return nil, errors.New("Invalid datatype")
		}
	}

	return Null{}, nil
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

// DeserializeInteger returns an Integer with value stored in it
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

// DeserializeBulkString returns a BulkString with value stored in it
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

// DeserializeArray returns an Array with value stored in it
func DeserializeArray(commands []string) (Array, int, error) {

	if commands[0][0] != '*' {
		return Array{}, 0, errors.New("Not an Array datatype")
	}
	val := commands[0][1:]
	arrLen, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("error:", err)
		return Array{}, 0, err
	}

	if arrLen == 0 {
		return Array{Length: 0, Items: []RESPType{}}, 1, nil
	}

	data := make([]RESPType, arrLen)
	elements := 0
	i := 0

	for i = 0; i < arrLen; i++ {
		cmd := addCRLF(commands[i+1:])
		v, el, err := deserializeCommand(cmd)
		if err != nil {
			return Array{}, 0, err
		}
		data[i] = v
		i += el
		elements++
	}

	if arrLen == elements && len(commands[i+1:]) == 0 {
		return Array{Length: arrLen, Items: data}, arrLen + 1, nil
	}

	return Array{}, 0, errors.New("Invalid Array: length mismatch")
}

func addCRLF(commands []string) string {
	str := ""

	for _, v := range commands {
		str += v + "\r\n"
	}

	return str
}

func deserializeCommand(cmds string) (RESPType, int, error) {
	// Empty command
	command := strings.Trim(cmds, " ")
	if len(command) == 0 {
		return nil, 0, errors.New("Empty input")
	}

	// No CRLF
	if !strings.HasSuffix(command, "\r\n") {
		return nil, 0, errors.New("No CRLF")
	}

	// Remove the last "\r\n" to remove an extra "" in the inputs array
	command = strings.TrimSuffix(command, "\r\n")

	// No commands
	inputs := strings.Split(command, "\r\n")
	if len(inputs) == 0 {
		return nil, 0, errors.New("Invalid input")
	}

	for i := 0; i < len(inputs); i++ {
		v := inputs[i]

		switch v[0] {
		case '+':
			val, err := DeserializeSimpleString(v)
			if err != nil {
				return nil, 0, err
			}
			return val, 0, nil
		case '-':
			val, err := DeserializeSimpleError(v)
			if err != nil {
				return nil, 0, err
			}
			return val, 0, nil
		case ':':
			val, err := DeserializeInteger(v)
			if err != nil {
				return nil, 0, err
			}
			return val, 0, nil
		case '$':
			val, el, err := DeserializeBulkString(inputs[i:])
			if err != nil {
				return nil, 0, err
			}
			return val, el, nil
		case '*':
			val, el, err := DeserializeArray(inputs[i:])
			if err != nil {
				return nil, 0, err
			}
			return val, el, nil
		default:
			return nil, 0, errors.New("Invalid datatype")
		}
	}

	return nil, 0, errors.New("Internal Error")
}
