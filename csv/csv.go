package csv

import (
	"encoding/csv"
	"os"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"writer": &interop.AdvFunction{
			Name:    "writer",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("file")},
			Value:   csvWriter,
		},
		"reader": &interop.AdvFunction{
			Name:    "reader",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("file")},
			Value:   csvReader,
		},
	}
}

func csvWriter(args interop.ArgMap) (tengo.Object, error) {
	file, _ := args.GetString("file")

	f, err := os.Create(file)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return makeCSVWriter(csv.NewWriter(f)), nil
}

func csvReader(args interop.ArgMap) (tengo.Object, error) {
	file, _ := args.GetString("file")

	f, err := os.Open(file)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return makeCSVReader(csv.NewReader(f)), nil
}
