//TODO: (1) Understand and use string.Builder instead of Sprintf

package resp

import (
	"errors"
	"fmt"
	"strings"
)

// RSPEType represents an interface that is extended by all other RSPE data types
type RSPEType interface {
	Serialize() (string, error)
}

//* Implementation of Simple Strings *//

// SimpleString represents an RESP simple string
type SimpleString struct {
	Value string
}

// Serialize returns the RESP serialization of the SimpleString, and an error
func (s SimpleString) Serialize() (string, error) {
	s.Value = strings.TrimSpace(s.Value)

	if strings.Contains(s.Value, "\r") || strings.Contains(s.Value, "\n") {
		return "", errors.New("SimpleString cannot contain \\r or \\n characters")
	} else if s.Value == "" {
		return "", errors.New("SimpleString cannot be empty")
	}

	return fmt.Sprintf("+%s\r\n", s.Value), nil
}

//* Implementation of Simple Errors *//

// SimpleError represents a RESP simple error
type SimpleError struct {
	Value string
}

// Serialize returns the RESP serialization of the SimpleError, and an error
func (s SimpleError) Serialize() (string, error) {
	s.Value = strings.TrimSpace(s.Value)

	if strings.Contains(s.Value, "\r") || strings.Contains(s.Value, "\n") {
		return "", errors.New("SimpleError cannot contain \\r or \\n characters")
	} else if s.Value == "" {
		return "", errors.New("SimpleError cannot be empty")
	}

	return fmt.Sprintf("-%s\r\n", s.Value), nil
}

//* Implementation of Integers *//

// Integer represents an RESP integer
type Integer struct {
	Value int
}

// Serialize returns the RESP serialization of the Integer, and an error
func (i Integer) Serialize() (string, error) {
	return fmt.Sprintf(":%d\r\n", i.Value), nil
}

//* Implementation of Bulk Strings *//

// BulkString returns a simple string in RESP
type BulkString struct {
	Value  string
	Length int
}

// Serialize returns the RESP serialization of the BulkString, and an error
func (bs BulkString) Serialize() (string, error) {
	bs.Value = strings.TrimSpace(bs.Value)
	length := len(bs.Value)
	bs.Length = length

	if bs.Length == 0 {
		return "", errors.New("BulkString cannot be empty")
	}

	return fmt.Sprintf("$%d\r\n%s\r\n", length, bs.Value), nil
}

//* Implementation of Arrays *//

// Array returns an array in RESP
type Array struct {
	Length int
	Items  []RSPEType
}

// Serialize returns the RESP serialization of the array, and an error
func (a Array) Serialize() (string, error) {
	length := len(a.Items)
	a.Length = length

	str := fmt.Sprintf("*%d\r\n", a.Length)
	for _, v := range a.Items {
		sv, err := v.Serialize()
		if err != nil {
			return "", err
		}

		str += sv
	}

	return str, nil
}
