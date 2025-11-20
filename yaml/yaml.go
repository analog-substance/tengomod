package yaml

import (
	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
	"gopkg.in/yaml.v3"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"encode": &interop.AdvFunction{
			Name:  "encode",
			Args:  []interop.AdvArg{interop.ObjectArg("input")},
			Value: encode,
		},
		"decode": &interop.AdvFunction{
			Name:  "decode",
			Args:  []interop.AdvArg{interop.UnionArg("input", interop.StrType, interop.ByteSliceType)},
			Value: decode,
		},
	}
}

func encode(args interop.ArgMap) (tengo.Object, error) {
	inputObj, _ := args.GetObject("input")
	input := tengo.ToInterface(inputObj)

	bytes, err := yaml.Marshal(&input)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return interop.GoByteSliceToTByteSlice(bytes)
}

func decode(args interop.ArgMap) (tengo.Object, error) {
	var input []byte

	inputStr, ok := args.GetString("input")
	if ok {
		input = []byte(inputStr)
	} else {
		input, _ = args.GetByteSlice("input")
	}

	output := make(map[string]interface{})
	err := yaml.Unmarshal(input, output)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	obj, err := interop.GoStrMapInterfaceToTMap(output)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}
	return obj, nil
}
