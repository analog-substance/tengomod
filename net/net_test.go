package net_test

import (
	"testing"

	"github.com/analog-substance/tengomod/internal/test"
)

func TestNet(t *testing.T) {
	test.Module(t, "net").Call("is_ip", "127.0.0.1").Expect(true)
	test.Module(t, "net").Call("is_ip", "example.com").Expect(false)
}
