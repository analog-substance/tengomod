package interop

import (
	"fmt"

	"github.com/analog-substance/tengo/v2"
)

func GoTSliceToGoInterfaceSlice(items []tengo.Object) []interface{} {
	var slice []interface{}
	for _, v := range items {
		slice = append(slice, v)
	}
	return slice
}

// GoStringSliceToTArray converts a golang string slice into a tengo Array
func GoStringSliceToTArray(slice []string) tengo.Object {
	var values []tengo.Object
	for _, s := range slice {
		values = append(values, &tengo.String{Value: s})
	}

	return &tengo.Array{
		Value: values,
	}
}

// GoIntSliceToTArray converts a golang int slice into a tengo Array
func GoIntSliceToTArray(slice []int) tengo.Object {
	var values []tengo.Object
	for _, i := range slice {
		values = append(values, &tengo.Int{Value: int64(i)})
	}

	return &tengo.Array{
		Value: values,
	}
}

// TArrayToGoStringSlice converts a tengo Array into a golang string slice
func TArrayToGoStringSlice(obj tengo.Object, name string) ([]string, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return GoTSliceToGoStringSlice(array.Value, name)
	case *tengo.ImmutableArray:
		return GoTSliceToGoStringSlice(array.Value, name)
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// TArrayToGoStringSlice converts a tengo Array into a golang string slice
func TArrayToGoInterfaceSlice(obj tengo.Object, name string) ([]interface{}, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return GoTSliceToGoInterfaceSlice(array.Value), nil
	case *tengo.ImmutableArray:
		return GoTSliceToGoInterfaceSlice(array.Value), nil
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// TArrayToGoSlice converts a tengo Array into a golang interface slice
func TArrayToGoSlice(obj tengo.Object, name string) ([]interface{}, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return GoTSliceToGoInterfaceSlice(array.Value), nil
	case *tengo.ImmutableArray:
		return GoTSliceToGoInterfaceSlice(array.Value), nil
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// TArrayToGoIntSlice converts a tengo Array into a golang int slice
func TArrayToGoIntSlice(obj tengo.Object, name string) ([]int, error) {
	switch array := obj.(type) {
	case *tengo.Array:
		return GoTSliceToGoIntSlice(array.Value, name)
	case *tengo.ImmutableArray:
		return GoTSliceToGoIntSlice(array.Value, name)
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "array(compatible)",
			Found:    obj.TypeName(),
		}
	}
}

// GoTSliceToGoStringSlice converts a slice of tengo Objects into a golang string slice
func GoTSliceToGoStringSlice(slice []tengo.Object, name string) ([]string, error) {
	var strSlice []string
	for idx, obj := range slice {
		item, ok := tengo.ToString(obj)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("%s[%d]", name, idx),
				Expected: "string(compatible)",
				Found:    obj.TypeName(),
			}
		}
		strSlice = append(strSlice, item)
	}
	return strSlice, nil
}

// GoTSliceToGoIntSlice converts a slice of tengo Objects into a golang int slice
func GoTSliceToGoIntSlice(slice []tengo.Object, name string) ([]int, error) {
	var intSlice []int
	for idx, obj := range slice {
		i, ok := tengo.ToInt(obj)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("%ss[%d]", name, idx),
				Expected: "int(compatible)",
				Found:    obj.TypeName(),
			}
		}

		intSlice = append(intSlice, i)
	}
	return intSlice, nil
}

// TMapToGoStringMapString converts a tengo object into a golang map[string]string
func TMapToGoStringMapString(obj tengo.Object, name string) (map[string]string, error) {
	var objMap map[string]tengo.Object
	switch o := obj.(type) {
	case *tengo.Map:
		objMap = o.Value
	case *tengo.ImmutableMap:
		objMap = o.Value
	default:
		return nil, tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "map(compatible)",
			Found:    obj.TypeName(),
		}
	}

	m := make(map[string]string)
	for key, value := range objMap {
		str, ok := tengo.ToString(value)
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     fmt.Sprintf("%s key %s", name, key),
				Expected: "string(compatible)",
				Found:    value.TypeName(),
			}
		}

		m[key] = str
	}

	return m, nil
}

// TStringToGoString converts a tengo object into a golang string
func TStringToGoString(arg tengo.Object, name string) (string, error) {
	str, ok := tengo.ToString(arg)
	if !ok {
		return "", tengo.ErrInvalidArgumentType{
			Name:     name,
			Expected: "string(compatible)",
			Found:    arg.TypeName(),
		}
	}

	return str, nil
}

// GoErrToTErr converts a golang error into a tengo Error
func GoErrToTErr(err error) tengo.Object {
	return &tengo.Error{
		Value: &tengo.String{
			Value: err.Error(),
		},
	}
}

func GoStringToTString(str string) tengo.Object {
	return &tengo.String{
		Value: str,
	}
}

type ArgValidation func([]tengo.Object) error

func WithExactArgs(n int) ArgValidation {
	return func(args []tengo.Object) error {
		if len(args) != n {
			return tengo.ErrWrongNumArguments
		}
		return nil
	}
}

func WithMinArgs(min int) ArgValidation {
	return func(args []tengo.Object) error {
		if len(args) < min {
			return tengo.ErrWrongNumArguments
		}
		return nil
	}
}

func WithMaxArgs(max int) ArgValidation {
	return func(args []tengo.Object) error {
		if len(args) > max {
			return tengo.ErrWrongNumArguments
		}
		return nil
	}
}

func WithArgRange(min int, max int) ArgValidation {
	return func(args []tengo.Object) error {
		if len(args) < min || len(args) > max {
			return tengo.ErrWrongNumArguments
		}
		return nil
	}
}

func NewCallable(callable tengo.CallableFunc, validations ...ArgValidation) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		for _, validation := range validations {
			err := validation(args)
			if err != nil {
				return nil, err
			}
		}
		return callable(args...)
	}
}

// FuncASSSSRSp transform a function of 'func(string, string, string, string) *string' signature
// into tengo CallableFunc type.
func FuncASSSSRSp(fn func(string, string, string, string) *string) tengo.CallableFunc {
	callable := func(args ...tengo.Object) (tengo.Object, error) {
		s1, err := TStringToGoString(args[0], "first")
		if err != nil {
			return nil, err
		}

		s2, err := TStringToGoString(args[1], "second")
		if err != nil {
			return nil, err
		}

		s3, err := TStringToGoString(args[2], "third")
		if err != nil {
			return nil, err
		}

		s4, err := TStringToGoString(args[3], "fourth")
		if err != nil {
			return nil, err
		}

		s := fn(s1, s2, s3, s4)
		if len(*s) > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		return &tengo.String{Value: *s}, nil
	}
	return NewCallable(callable, WithExactArgs(4))
}

// FuncASSSRSp transform a function of 'func(string, string, string) *string' signature
// into tengo CallableFunc type.
func FuncASSSRSp(fn func(string, string, string) *string) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 3 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		s2, ok := tengo.ToString(args[1])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "string(compatible)",
				Found:    args[1].TypeName(),
			}
		}
		s3, ok := tengo.ToString(args[2])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "third",
				Expected: "string(compatible)",
				Found:    args[2].TypeName(),
			}
		}
		s := fn(s1, s2, s3)
		if len(*s) > tengo.MaxStringLen {
			return nil, tengo.ErrStringLimit
		}
		return &tengo.String{Value: *s}, nil
	}
}

// FuncASSSsSRSsp transform a function of 'func(string, string, []string, string) *[]string' signature
// into tengo CallableFunc type.
func FuncASSSsSRSsp(fn func(string, string, []string, string) *[]string) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 4 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		s2, ok := tengo.ToString(args[1])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "string(compatible)",
				Found:    args[1].TypeName(),
			}
		}

		var ss1 []string
		switch arg2 := args[2].(type) {
		case *tengo.Array:
			for idx, a := range arg2.Value {
				as, ok := tengo.ToString(a)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{
						Name:     fmt.Sprintf("third[%d]", idx),
						Expected: "string(compatible)",
						Found:    a.TypeName(),
					}
				}
				ss1 = append(ss1, as)
			}
		case *tengo.ImmutableArray:
			for idx, a := range arg2.Value {
				as, ok := tengo.ToString(a)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{
						Name:     fmt.Sprintf("third[%d]", idx),
						Expected: "string(compatible)",
						Found:    a.TypeName(),
					}
				}
				ss1 = append(ss1, as)
			}
		default:
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "third",
				Expected: "array",
				Found:    args[0].TypeName(),
			}
		}

		s4, ok := tengo.ToString(args[3])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "fourth",
				Expected: "string(compatible)",
				Found:    args[3].TypeName(),
			}
		}

		arr := &tengo.Array{}
		for _, res := range *fn(s1, s2, ss1, s4) {
			if len(res) > tengo.MaxStringLen {
				return nil, tengo.ErrStringLimit
			}
			arr.Value = append(arr.Value, &tengo.String{Value: res})
		}
		return arr, nil
	}
}

// FuncASSsSRSsp transform a function of 'func(string, []string, string) *[]string' signature
// into tengo CallableFunc type.
func FuncASSsSRSsp(fn func(string, []string, string) *[]string) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 3 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		var ss1 []string
		switch arg1 := args[1].(type) {
		case *tengo.Array:
			for idx, a := range arg1.Value {
				as, ok := tengo.ToString(a)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{
						Name:     fmt.Sprintf("second[%d]", idx),
						Expected: "string(compatible)",
						Found:    a.TypeName(),
					}
				}
				ss1 = append(ss1, as)
			}
		case *tengo.ImmutableArray:
			for idx, a := range arg1.Value {
				as, ok := tengo.ToString(a)
				if !ok {
					return nil, tengo.ErrInvalidArgumentType{
						Name:     fmt.Sprintf("second[%d]", idx),
						Expected: "string(compatible)",
						Found:    a.TypeName(),
					}
				}
				ss1 = append(ss1, as)
			}
		default:
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "array",
				Found:    args[1].TypeName(),
			}
		}

		s3, ok := tengo.ToString(args[2])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "third",
				Expected: "string(compatible)",
				Found:    args[2].TypeName(),
			}
		}

		arr := &tengo.Array{}
		for _, res := range *fn(s1, ss1, s3) {
			if len(res) > tengo.MaxStringLen {
				return nil, tengo.ErrStringLimit
			}
			arr.Value = append(arr.Value, &tengo.String{Value: res})
		}
		return arr, nil
	}
}

// FuncASSBSRBp transform a function of 'func(string, string, bool, string) *bool' signature
// into tengo CallableFunc type.
func FuncASSBSRBp(fn func(string, string, bool, string) *bool) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 4 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		s2, ok := tengo.ToString(args[1])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "string(compatible)",
				Found:    args[1].TypeName(),
			}
		}

		b1, ok := tengo.ToBool(args[2])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "third",
				Expected: "bool(compatible)",
				Found:    args[2].TypeName(),
			}
		}

		s4, ok := tengo.ToString(args[3])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "fourth",
				Expected: "string(compatible)",
				Found:    args[3].TypeName(),
			}
		}

		if *fn(s1, s2, b1, s4) {
			return tengo.TrueValue, nil
		}
		return tengo.FalseValue, nil
	}
}

// FuncASBSRBp transform a function of 'func(string, bool, string) *bool' signature
// into tengo CallableFunc type.
func FuncASBSRBp(fn func(string, bool, string) *bool) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 3 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		b1, ok := tengo.ToBool(args[1])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "bool(compatible)",
				Found:    args[1].TypeName(),
			}
		}

		s4, ok := tengo.ToString(args[2])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "third",
				Expected: "string(compatible)",
				Found:    args[2].TypeName(),
			}
		}

		if *fn(s1, b1, s4) {
			return tengo.TrueValue, nil
		}
		return tengo.FalseValue, nil
	}
}

// FuncASSISRIp transform a function of 'func(string, string, int, string) *int' signature
// into tengo CallableFunc type.
func FuncASSISRIp(fn func(string, string, int, string) *int) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 4 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		s2, ok := tengo.ToString(args[1])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "string(compatible)",
				Found:    args[1].TypeName(),
			}
		}

		i1, ok := tengo.ToInt(args[2])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "third",
				Expected: "int(compatible)",
				Found:    args[2].TypeName(),
			}
		}

		s4, ok := tengo.ToString(args[3])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "fourth",
				Expected: "string(compatible)",
				Found:    args[3].TypeName(),
			}
		}

		i := fn(s1, s2, i1, s4)
		return &tengo.Int{Value: int64(*i)}, nil
	}
}

// FuncASISRIp transform a function of 'func(string, int, string) *int' signature
// into tengo CallableFunc type.
func FuncASISRIp(fn func(string, int, string) *int) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 3 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		i1, ok := tengo.ToInt(args[1])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "second",
				Expected: "int(compatible)",
				Found:    args[1].TypeName(),
			}
		}

		s4, ok := tengo.ToString(args[2])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "third",
				Expected: "string(compatible)",
				Found:    args[2].TypeName(),
			}
		}

		i := fn(s1, i1, s4)
		return &tengo.Int{Value: int64(*i)}, nil
	}
}

// FuncASRSsE transform a function of 'func(string) ([]string, error)' signature
// into tengo CallableFunc type.
func FuncASRSsE(fn func(string) ([]string, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		res, err := fn(s1)
		if err != nil {
			return GoErrToTErr(err), nil
		}

		arr := &tengo.Array{}
		for _, r := range res {
			if len(r) > tengo.MaxStringLen {
				return nil, tengo.ErrStringLimit
			}
			arr.Value = append(arr.Value, &tengo.String{Value: r})
		}
		return arr, nil
	}
}

// FuncASRBE transform a function of 'func(string) (bool, error)' signature
// into tengo CallableFunc type.
func FuncASRBE(fn func(string) (bool, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		res, err := fn(s1)
		if err != nil {
			return GoErrToTErr(err), nil
		}

		if res {
			return tengo.TrueValue, nil
		}
		return tengo.FalseValue, nil
	}
}

// FuncASRB transform a function of 'func(string) bool' signature
// into tengo CallableFunc type.
func FuncASRB(fn func(string) bool) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}

		res := fn(s1)
		if res {
			return tengo.TrueValue, nil
		}
		return tengo.FalseValue, nil
	}
}

// FuncASvRSsE transform a function of 'func(...string) ([]string, error)' signature
// into tengo CallableFunc type.
func FuncASvRSsE(fn func(...string) ([]string, error)) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) == 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		var strings []string
		for i, arg := range args {
			str, ok := tengo.ToString(arg)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     fmt.Sprintf("#%d arg", i),
					Expected: "string(compatible)",
					Found:    arg.TypeName(),
				}
			}

			strings = append(strings, str)
		}

		res, err := fn(strings...)
		if err != nil {
			return GoErrToTErr(err), nil
		}

		return GoStringSliceToTArray(res), nil
	}
}

// FuncASvRB transform a function of 'func(...string) bool' signature
// into tengo CallableFunc type.
func FuncASvRB(fn func(...string) bool) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) == 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		var strings []string
		for i, arg := range args {
			str, ok := tengo.ToString(arg)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     fmt.Sprintf("#%d arg", i),
					Expected: "string(compatible)",
					Found:    arg.TypeName(),
				}
			}

			strings = append(strings, str)
		}

		res := fn(strings...)
		if res {
			return tengo.TrueValue, nil
		}

		return tengo.FalseValue, nil
	}
}

// FuncASvRS transform a function of 'func(...string) string' signature
// into tengo CallableFunc type.
func FuncASvRS(fn func(...string) string) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) == 0 {
			return nil, tengo.ErrWrongNumArguments
		}
		var strings []string
		for i, arg := range args {
			str, ok := tengo.ToString(arg)
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     fmt.Sprintf("#%d arg", i),
					Expected: "string(compatible)",
					Found:    arg.TypeName(),
				}
			}

			strings = append(strings, str)
		}

		return &tengo.String{Value: fn(strings...)}, nil
	}
}

// FuncASRI transform a function of 'func(string) int' signature into
// CallableFunc type.
func FuncASRI(fn func(string) int) tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		if len(args) != 1 {
			return nil, tengo.ErrWrongNumArguments
		}
		s1, ok := tengo.ToString(args[0])
		if !ok {
			return nil, tengo.ErrInvalidArgumentType{
				Name:     "first",
				Expected: "string(compatible)",
				Found:    args[0].TypeName(),
			}
		}
		i := fn(s1)
		return &tengo.Int{Value: int64(i)}, nil
	}
}
