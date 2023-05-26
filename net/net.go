package net

import (
	"net"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"is_ip": &interop.AdvFunction{
			Name:    "is_ip",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("input")},
			Value:   isIP,
		},
	}
}

func isIP(args interop.ArgMap) (tengo.Object, error) {
	input, _ := args.GetString("input")
	parsed := net.ParseIP(input)
	return interop.GoBoolToTBool(parsed == nil), nil
}
