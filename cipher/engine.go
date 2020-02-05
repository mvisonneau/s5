package cipher

import (
	"fmt"
	"regexp"
)

const (
	InputRegexp string = `{{\s?s5:([A-Za-z0-9+\\/=]*)\s?}}`
)

// Engine is an interface of supported/required commands for each cipher engine
type Engine interface {
	Cipher(string) (string, error)
	Decipher(string) (string, error)
}

func GenerateOutput(value string) string {
	return fmt.Sprintf("{{s5:%s}}", value)
}

func ParseInput(value string) (string, error) {
	re := regexp.MustCompile(InputRegexp)
	if !re.MatchString(value) {
		return "", fmt.Errorf("Invalid string format, should be '{{s5:*}}'")
	}
	return re.FindStringSubmatch(value)[1], nil
}
