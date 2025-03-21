package resp

import "strings"

// GetCommands extracts all commands from the input string and returns them as a slice of strings
func GetCommands(buf []byte) []string {
	val := string(buf)
	val = strings.ToUpper(val)

	inputs := strings.Split(val, "\r\n")

	return inputs
}
