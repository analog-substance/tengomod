package csv

import (
	"encoding/csv"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengo/v2/stdlib"
	"github.com/analog-substance/tengomod/interop"
	"github.com/analog-substance/tengomod/types"
)

type CSVWriter struct {
	types.PropObject
	Value *csv.Writer
}

// TypeName should return the name of the type.
func (w *CSVWriter) TypeName() string {
	return "csv-writer"
}

// String should return a string representation of the type's value.
func (w *CSVWriter) String() string {
	return "<csv-writer>"
}

// IsFalsy should return true if the value of the type should be considered
// as falsy.
func (w *CSVWriter) IsFalsy() bool {
	return w.Value == nil
}

// CanIterate should return whether the Object can be Iterated.
func (w *CSVWriter) CanIterate() bool {
	return false
}

func (w *CSVWriter) write(args interop.ArgMap) (tengo.Object, error) {
	row, _ := args.GetStringSlice("row")

	err := w.Value.Write(row)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

func (w *CSVWriter) writeAll(args interop.ArgMap) (tengo.Object, error) {
	rows, _ := args.GetStringSliceSlice("rows")

	err := w.Value.WriteAll(rows)
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return nil, nil
}

func makeCSVWriter(w *csv.Writer) *CSVWriter {
	writer := &CSVWriter{
		Value: w,
	}

	objectMap := map[string]tengo.Object{
		"write": &interop.AdvFunction{
			Name:    "write",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrSliceArg("row", false)},
			Value:   writer.write,
		},
		"write_all": &interop.AdvFunction{
			Name:    "write_all",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrSliceArg("rows", false)},
			Value:   writer.writeAll,
		},
		"flush": &tengo.UserFunction{
			Name:  "flush",
			Value: stdlib.FuncAR(writer.Value.Flush),
		},
	}

	writer.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: make(map[string]types.Property),
	}

	return writer
}
