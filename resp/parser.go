//TODO: (1) Change the return type to be based on the input, like if the input is -me\r\n don't return an array

package resp

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Deserialize takes a command as input and deserializes it
func Deserialize(cmds string) (Type, error) {
	// Empty command
	command := strings.Trim(cmds, " ")
	if len(command) == 0 {
		return nil, errors.New("empty input")
	}

	// No CRLF
	if !strings.HasSuffix(command, "\r\n") {
		return nil, errors.New("no CRLF")
	}

	// Remove the last "\r\n" to remove an extra "" in the inputs array
	command = strings.TrimSuffix(command, "\r\n")

	// No commands
	inputs := strings.Split(command, "\r\n")
	if len(inputs) == 0 {
		return nil, errors.New("invalid input")
	}

	for i := range len(inputs) {
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
			return nil, errors.New("invalid datatype")
		}
	}

	return Null{}, nil
}

// DeserializeSimpleString returns a simpleString with value stored in it
func DeserializeSimpleString(command string) (SimpleString, error) {
	if command[0] != '+' {
		return SimpleString{}, errors.New("not SimpleString datatype")
	}
	s := SimpleString{Value: command[1:]}
	return s, nil
}

// DeserializeSimpleError returns a simpleError with value stored in it
func DeserializeSimpleError(command string) (SimpleError, error) {
	if command[0] != '-' {
		return SimpleError{}, errors.New("not SimpleError datatype")
	}
	s := SimpleError{Value: command[1:]}
	return s, nil
}

// DeserializeInteger returns an Integer with value stored in it
func DeserializeInteger(command string) (Integer, error) {
	if command[0] != ':' {
		return Integer{}, errors.New("not Integer datatype")
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
		return BulkString{}, 0, errors.New("not a BulkString datatype")
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
			return BulkString{}, 0, errors.New("invalid Bulk String: length mismatch")
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
	if len(commands) == 0 || commands[0][0] != '*' {
		return Array{}, 0, errors.New("not an Array datatype")
	}

	// Extract array length
	val := commands[0][1:]
	arrLen, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("error:", err)
		return Array{}, 0, err
	}

	// Handle empty array
	if arrLen == 0 {
		return Array{Length: 0, Items: []Type{}}, 1, nil
	}

	data := make([]Type, arrLen)
	elements := 0
	index := 0 // Handling the indexes in CMDS array

	for i := range arrLen {
		cmd := addCRLF(commands[index+1:])
		v, el, err := deserializeCommand(cmd)
		if err != nil {
			return Array{}, 0, err
		}

		data[i] = v
		index += el + 1
		elements++
	}

	if arrLen == elements && len(commands[index+1:]) == 0 {
		return Array{Length: arrLen, Items: data}, arrLen + 1, nil
	}

	return Array{}, 0, errors.New("invalid Array: length mismatch")
}

func addCRLF(commands []string) string {
	str := ""

	for _, v := range commands {
		str += v + "\r\n"
	}

	return str
}

func deserializeCommand(cmds string) (Type, int, error) {
	// Empty command
	command := strings.Trim(cmds, " ")
	if len(command) == 0 {
		return nil, 0, errors.New("empty input")
	}

	// No CRLF
	if !strings.HasSuffix(command, "\r\n") {
		return nil, 0, errors.New("no CRLF")
	}

	// Remove the last "\r\n" to remove an extra "" in the inputs array
	command = strings.TrimSuffix(command, "\r\n")

	// No commands
	inputs := strings.Split(command, "\r\n")
	if len(inputs) == 0 {
		return nil, 0, errors.New("invalid input")
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
			fmt.Println("v[0] :", v[0])
			return nil, 0, errors.New("invalid datatype")
		}
	}

	return nil, 0, errors.New("internal Error")
}
