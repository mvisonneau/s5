package main

import (
	"testing"
)

func TestRunCli(t *testing.T) {
	c := runCli()
	if c.Name != "s5" {
		t.Fatalf("Expected c.Name to be s5, got '%v'", c.Name)
	}
}
