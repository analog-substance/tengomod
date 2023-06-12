package csv

import (
	"encoding/csv"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
	"github.com/analog-substance/tengomod/types"
)

type CSVReader struct {
	types.PropObject
	Value *csv.Reader
}

// TypeName should return the name of the type.
func (r *CSVReader) TypeName() string {
	return "csv-reader"
}

// String should return a string representation of the type's value.
func (r *CSVReader) String() string {
	return "<csv-reader>"
}

// IsFalsy should return true if the value of the type should be considered
// as falsy.
func (r *CSVReader) IsFalsy() bool {
	return r.Value == nil
}

// CanIterate should return whether the Object can be Iterated.
func (r *CSVReader) CanIterate() bool {
	return false
}

func (r *CSVReader) read(args interop.ArgMap) (tengo.Object, error) {
	row, err := r.Value.Read()
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return interop.GoStrSliceToTArray(row), nil
}

func (w *CSVReader) readAll(args interop.ArgMap) (tengo.Object, error) {
	rows, err := w.Value.ReadAll()
	if err != nil {
		return interop.GoErrToTErr(err), nil
	}

	return interop.GoStrSliceSliceToTArray(rows), nil
}

func makeCSVReader(r *csv.Reader) *CSVReader {
	reader := &CSVReader{
		Value: r,
	}

	objectMap := map[string]tengo.Object{
		"read": &interop.AdvFunction{
			Name:  "read",
			Value: reader.read,
		},
		"read_all": &interop.AdvFunction{
			Name:  "read_all",
			Value: reader.readAll,
		},
	}

	reader.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: make(map[string]types.Property),
	}

	return reader
}
