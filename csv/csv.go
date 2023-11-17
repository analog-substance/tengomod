package csv

import (
	"encoding/csv"
	"os"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
)

func Module() map[string]tengo.Object {
	return map[string]tengo.Object{
		"write": &interop.AdvFunction{
			Name:    "write",
			NumArgs: interop.ExactArgs(2),
			Args: []interop.AdvArg{
				interop.StrArg("file"),
				interop.UnionArg("data", interop.StrSliceType, interop.StrSliceSliceType),
			},
			Value: csvWrite,
		},
		"writer": &interop.AdvFunction{
			Name:    "writer",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("file")},
			Value:   csvWriter,
		},
		"read": &interop.AdvFunction{
			Name:    "read",
			NumArgs: interop.ExactArgs(1),
			Args: []interop.AdvArg{
				interop.StrArg("file"),
			},
			Value: csvRead,
		},
		"reader": &interop.AdvFunction{
			Name:    "reader",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("file")},
			Value:   csvReader,
		},
	}
}

func csvWrite(args interop.ArgMap) (tengo.Object, error) {
	file, _ := args.GetString("file")

	f, err := os.Create(file)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	writer := makeCSVWriter(csv.NewWriter(f))

	rows, ok := args.GetStringSliceSlice("data")
	if !ok {
		row, _ := args.GetStringSlice("data")
		rows = append(rows, row)
	}

	return writer.writeAll(interop.ArgMap{
		"rows": rows,
	})
}

func csvWriter(args interop.ArgMap) (tengo.Object, error) {
	file, _ := args.GetString("file")

	f, err := os.Create(file)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return makeCSVWriter(csv.NewWriter(f)), nil
}

func csvRead(args interop.ArgMap) (tengo.Object, error) {
	file, _ := args.GetString("file")

	f, err := os.Open(file)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	reader := makeCSVReader(csv.NewReader(f))

	return reader.readAll(make(interop.ArgMap))
}

func csvReader(args interop.ArgMap) (tengo.Object, error) {
	file, _ := args.GetString("file")

	f, err := os.Open(file)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return makeCSVReader(csv.NewReader(f)), nil
}
