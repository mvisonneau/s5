package main

var cfg Config

// Config : A type that handles the configuration of the app
type Config struct {
	Log struct {
		Level  string
		Format string
	}

	TransitKey string

	Vault struct {
		Address string
		Token   string
	}
}
