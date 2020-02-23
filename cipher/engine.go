package cipher

import (
	"fmt"
	"regexp"
)

const (
	// InputRegexp is defining the syntax of an s5 input variable
	InputRegexp string = `{{\s?s5:([A-Za-z0-9+\\/=]*)\s?}}`
)

// Engine is an interface of supported/required commands for each cipher engine
type Engine interface {
	Cipher(string) (string, error)
	Decipher(string) (string, error)
}

// GenerateOutput return a ciphered string in a s5 format
func GenerateOutput(value string) string {
	return fmt.Sprintf("{{s5:%s}}", value)
}

// ParseInput retrieves ciphered value from a string in the s5 format
func ParseInput(value string) (string, error) {
	re := regexp.MustCompile(InputRegexp)
	if !re.MatchString(value) {
		return "", fmt.Errorf("Invalid string format, should be '{{s5:*}}'")
	}
	return re.FindStringSubmatch(value)[1], nil
}
