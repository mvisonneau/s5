package cipher

// Engine is an interface of supported/required commands for each cipher engine
type Engine interface {
	Cipher(string) (string, error)
	Decipher(string) (string, error)
}
