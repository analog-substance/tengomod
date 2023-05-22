package viper

import (
	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/stdlib"
	"github.com/analog-substance/tengomod/interop"
	"github.com/spf13/viper"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"get_string": &tengo.UserFunction{
			Name:  "get_string",
			Value: stdlib.FuncASRS(viper.GetString),
		},
		"get_int": &tengo.UserFunction{
			Name:  "get_int",
			Value: interop.FuncASRI(viper.GetInt),
		},
		"get_bool": &tengo.UserFunction{
			Name:  "get_bool",
			Value: interop.FuncASRB(viper.GetBool),
		},
	}
}
