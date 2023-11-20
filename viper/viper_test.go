package viper_test

import (
	"testing"

	"github.com/analog-substance/tengomod/internal/test"
	"github.com/spf13/viper"
)

func TestViper(t *testing.T) {
	viper.Set("test.string", "tengomod")
	viper.Set("test.int", 123)
	viper.Set("test.bool", true)

	test.Module(t, "viper").Call("get_string", "test.string").Expect("tengomod")
	test.Module(t, "viper").Call("get_int", "test.int").Expect(123)
	test.Module(t, "viper").Call("get_bool", "test.bool").Expect(true)
}
