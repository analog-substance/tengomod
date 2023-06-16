package http

import (
	"net/http"

	"github.com/analog-substance/tengo/v2"
	"github.com/analog-substance/tengomod/interop"
	"github.com/analog-substance/tengomod/types"
)

type HTTPHeader struct {
	types.PropObject
	Value http.Header
}

func (r *HTTPHeader) TypeName() string {
	return "http-header"
}

// String should return a string representation of the type's value.
func (r *HTTPHeader) String() string {
	return "<http-header>"
}

// IsFalsy should return true if the value of the type should be considered
// as falsy.
func (h *HTTPHeader) IsFalsy() bool {
	return h.Value == nil
}

// CanIterate should return whether the Object can be Iterated.
func (h *HTTPHeader) CanIterate() bool {
	return true
}

// Iterate returns an iterator.
func (h *HTTPHeader) Iterate() tengo.Iterator {
	value := make(map[string]tengo.Object)
	for header, values := range h.Value {
		value[header] = interop.GoStrSliceToTArray(values)
	}
	m := &tengo.Map{
		Value: value,
	}

	return m.Iterate()
}

func (h *HTTPHeader) add(args interop.ArgMap) (tengo.Object, error) {
	key, _ := args.GetString("key")
	value, _ := args.GetString("value")

	h.Value.Add(key, value)
	return nil, nil
}

func (h *HTTPHeader) set(args interop.ArgMap) (tengo.Object, error) {
	key, _ := args.GetString("key")
	value, _ := args.GetString("value")

	h.Value.Set(key, value)
	return nil, nil
}

func (h *HTTPHeader) get(args interop.ArgMap) (tengo.Object, error) {
	key, _ := args.GetString("key")

	return interop.GoStrToTStr(h.Value.Get(key)), nil
}

func (h *HTTPHeader) delete(args interop.ArgMap) (tengo.Object, error) {
	key, _ := args.GetString("key")

	h.Value.Del(key)
	return nil, nil
}

func (h *HTTPHeader) values(args interop.ArgMap) (tengo.Object, error) {
	key, _ := args.GetString("key")

	return interop.GoStrSliceToTArray(h.Value.Values(key)), nil
}

func makeHTTPHeader(h http.Header) *HTTPHeader {
	header := &HTTPHeader{
		Value: h,
	}

	objectMap := map[string]tengo.Object{
		"add": &interop.AdvFunction{
			Name:    "add",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrArg("key"), interop.StrArg("value")},
			Value:   header.add,
		},
		"set": &interop.AdvFunction{
			Name:    "set",
			NumArgs: interop.ExactArgs(2),
			Args:    []interop.AdvArg{interop.StrArg("key"), interop.StrArg("value")},
			Value:   header.set,
		},
		"get": &interop.AdvFunction{
			Name:    "get",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("key")},
			Value:   header.get,
		},
		"del": &interop.AdvFunction{
			Name:    "del",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("key")},
			Value:   header.delete,
		},
		"values": &interop.AdvFunction{
			Name:    "values",
			NumArgs: interop.ExactArgs(1),
			Args:    []interop.AdvArg{interop.StrArg("key")},
			Value:   header.values,
		},
	}

	header.PropObject = types.PropObject{
		ObjectMap:  objectMap,
		Properties: make(map[string]types.Property),
	}

	return header
}
