//TODO: (1) Understand and use string.Builder instead of Sprintf

package resp

import (
	"errors"
	"fmt"
	"strings"
)

// RSPEType represents an interface that is extended by all other data types
type RSPEType interface {
	Serialize() (string, error)
}

//* Implementation of Simple Strings *//

// SimpleString returns a simple string in RESP
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

// SimpleError returns a simple error in RESP
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

// Integer returns a simple error in RESP
type Integer struct {
	Value int
}

// Serialize returns the RESP serialization of the SimpleError, and an error
func (i Integer) Serialize() (string, error) {
	return fmt.Sprintf(":%d\r\n", i.Value), nil
}
