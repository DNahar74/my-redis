package resp

import (
	"errors"
	"strconv"
	"strings"
)

// Type represents an interface that is extended by all other RESP data types
type Type interface {
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
		return "", errors.New("SimpleString cannot contain CR or LF characters")
	} else if s.Value == "" {
		return "", errors.New("SimpleString cannot be empty")
	}

	var sb strings.Builder
	sb.WriteByte('+')
	sb.WriteString(s.Value)
	sb.WriteString("\r\n")
	return sb.String(), nil
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

	var sb strings.Builder
	sb.WriteByte('-')
	sb.WriteString(s.Value)
	sb.WriteString("\r\n")
	return sb.String(), nil
}

//* Implementation of Integers *//

// Integer represents an RESP integer
type Integer struct {
	Value int
}

// Serialize returns the RESP serialization of the Integer, and an error
func (i Integer) Serialize() (string, error) {
	var sb strings.Builder
	sb.WriteByte(':')
	sb.WriteString(strconv.Itoa(i.Value))
	sb.WriteString("\r\n")
	return sb.String(), nil
}

//* Implementation of Bulk Strings *//

// BulkString returns a simple string in RESP
type BulkString struct {
	Value  string
	Length int
}

// Serialize returns the RESP serialization of the BulkString, and an error
func (bs BulkString) Serialize() (string, error) {
	length := len(bs.Value)
	bs.Length = length

	var sb strings.Builder
	sb.WriteByte('$')
	sb.WriteString(strconv.Itoa(bs.Length))
	sb.WriteString("\r\n")
	sb.WriteString(bs.Value)
	sb.WriteString("\r\n")
	return sb.String(), nil
}

//* Implementation of Arrays *//

// Array returns an array in RESP
type Array struct {
	Length int
	Items  []Type
}

// Serialize returns the RESP serialization of the array, and an error
func (a Array) Serialize() (string, error) {
	length := len(a.Items)
	a.Length = length

	var sb strings.Builder
	sb.WriteByte('*')
	sb.WriteString(strconv.Itoa(a.Length))
	sb.WriteString("\r\n")

	for _, v := range a.Items {
		sv, err := v.Serialize()
		if err != nil {
			return "", err
		}

		sb.WriteString(sv)
	}

	return sb.String(), nil
}

//* Implementation of Null *//

// Null represents a RESP Null value
type Null struct{}

// Serialize returns the RESP serialization of the Null value, and an error
func (s Null) Serialize() (string, error) {
	return "_\r\n", nil
}
